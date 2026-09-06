package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"{{MODULE_NAME}}/internal/app"
	"{{MODULE_NAME}}/routes"

	"github.com/velocitykode/velocity"
	"github.com/velocitykode/velocity/auth"
	"github.com/velocitykode/velocity/auth/drivers/schemes"
	"github.com/velocitykode/velocity/router"
	velhttp "github.com/velocitykode/velocity/testing/http"
	"github.com/velocitykode/velocity/velocitytest"
)

// testPeer is the socket address httptest.NewRequest stamps on every
// request. The "trusted chain" case below lists it as the proxy.
const testPeer = "192.0.2.1"

// observedIdentity is what a handler sees for the client identity of a
// request after the full global and API middleware stacks ran.
type observedIdentity struct {
	ip         string
	host       string
	remoteAddr string
	scheme     string
}

// proxyTestApp boots the real application (modules, middleware stacks,
// routes) the way main.go does, minus the listener, from the same env
// configuration path production uses: AUTH_SCHEME=api builds the JWT
// scheme (the template's default) and AUTH_TRUSTED_PROXIES installs the
// proxy trust set in the router (Context.IP) and the auth manager (login
// throttle keys). The in-memory cache also installs the framework's
// default login throttler. Two probe routes under the API stack report
// the identity a handler observes and run the scheme's Attempt, the call
// any login endpoint added to this template makes.
func proxyTestApp(t *testing.T, trustedProxies []string) (*velocity.App, *observedIdentity) {
	t.Helper()
	t.Setenv("APP_ENV", "testing")
	t.Setenv("APP_KEY", "base64:"+base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x2a}, 32)))
	t.Setenv("AUTH_SCHEME", "api")
	t.Setenv("AUTH_JWT_SECRET", base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x3b}, 48)))
	t.Setenv("AUTH_JWT_BLACKLIST_ENABLED", "false") // the template's shipped default; no blacklist store is wired here
	t.Setenv("AUTH_TRUSTED_PROXIES", strings.Join(trustedProxies, ","))
	t.Setenv("DB_CONNECTION", "") // no database: the probes and the login flow never reach it
	t.Setenv("CACHE_DRIVER", "memory")
	t.Setenv("QUEUE_DRIVER", "memory")
	t.Setenv("LOG_DRIVER", "console")
	t.Setenv("MAIL_DRIVER", "log")
	a, err := velocitytest.NewApp(velocity.WithConfig(velocity.ConfigFromEnv()))
	if err != nil {
		t.Fatalf("velocity.New: %v", err)
	}
	t.Cleanup(func() { _ = a.Shutdown(context.Background()) })

	var seen observedIdentity
	a.Modules(app.Configure).
		Middleware(app.Middleware).
		Routes(func(r *velocity.Routing) {
			routes.Register(r)
			r.API("/_probe", func(api router.Router) {
				api.Get("/identity", func(c *router.Context) error {
					seen = observedIdentity{ip: c.IP(), host: c.Request.Host, remoteAddr: c.Request.RemoteAddr, scheme: c.Request.URL.Scheme}
					return c.JSON(http.StatusOK, map[string]string{"ip": seen.ip})
				})
				api.Post("/login", func(c *router.Context) error {
					var body struct {
						Email    string `json:"email"`
						Password string `json:"password"`
					}
					if err := c.Bind(&body); err != nil {
						return router.NewHTTPError(http.StatusBadRequest, "invalid body")
					}
					ok, err := auth.FromContext(c).Attempt(c.Response, c.Request, map[string]any{"email": body.Email, "password": body.Password})
					if err != nil || !ok {
						return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
					}
					return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
				})
			})
		})
	if err := a.Bootstrap(); err != nil {
		t.Fatalf("app.Bootstrap: %v", err)
	}
	return a, &seen
}

func probe(t *testing.T, a *velocity.App, host string, headers map[string]string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/_probe/identity", nil)
	req.Host = host
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	a.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("probe status = %d, want 200", rec.Code)
	}
}

// TestMiddlewareStack_UntrustedPeerCannotChangeIdentity runs spoofed
// forwarding headers through the application's real middleware stacks
// from a peer that is not a configured proxy: the resolved IP, the
// request host, the socket address and the URL scheme all stay as the
// socket delivered them. X-Forwarded-Proto in particular never reaches
// URL.Scheme.
func TestMiddlewareStack_UntrustedPeerCannotChangeIdentity(t *testing.T) {
	cases := []struct {
		name           string
		trustedProxies []string
	}{
		{name: "no proxies configured", trustedProxies: nil},
		{name: "peer outside the configured proxy set", trustedProxies: []string{"10.0.0.0/8"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			a, seen := proxyTestApp(t, tc.trustedProxies)
			probe(t, a, "app.example", map[string]string{
				"X-Forwarded-For":   "203.0.113.5",
				"X-Forwarded-Host":  "attacker.example",
				"X-Forwarded-Proto": "https",
				"Forwarded":         "for=203.0.113.6;host=attacker.example;proto=https",
			})
			if seen.ip != testPeer {
				t.Fatalf("IP = %q, want socket peer %s", seen.ip, testPeer)
			}
			if seen.host != "app.example" {
				t.Fatalf("Host = %q, want app.example", seen.host)
			}
			if seen.remoteAddr != testPeer+":1234" {
				t.Fatalf("RemoteAddr = %q, want untouched socket address", seen.remoteAddr)
			}
			if seen.scheme != "" {
				t.Fatalf("URL.Scheme = %q, want empty (never taken from X-Forwarded-Proto)", seen.scheme)
			}
		})
	}
}

// TestMiddlewareStack_TrustedChainResolvesFromTrustedEnd sends a request
// through a configured proxy with an attacker-supplied prefix ahead of
// the real client: the identity is the right-most hop that is not a
// trusted proxy, never the prefix, and the host and scheme are never
// rewritten.
func TestMiddlewareStack_TrustedChainResolvesFromTrustedEnd(t *testing.T) {
	a, seen := proxyTestApp(t, []string{testPeer})
	probe(t, a, "app.example", map[string]string{
		"X-Forwarded-For":   "203.0.113.9, 198.51.100.7, " + testPeer,
		"X-Forwarded-Host":  "attacker.example",
		"X-Forwarded-Proto": "https",
	})
	if seen.ip != "198.51.100.7" {
		t.Fatalf("IP = %q, want 198.51.100.7 (right-most untrusted hop, not the attacker prefix)", seen.ip)
	}
	if seen.host != "app.example" {
		t.Fatalf("Host = %q, want app.example (X-Forwarded-Host is never applied)", seen.host)
	}
	if seen.scheme != "" {
		t.Fatalf("URL.Scheme = %q, want empty (scheme is not taken from a forwarding header)", seen.scheme)
	}
}

// countingUsers is an auth.UserStore that knows one user and counts how
// many candidate passwords reached verification.
type countingUsers struct {
	mu     sync.Mutex
	checks int
}

type probeUser struct{}

func (*probeUser) GetAuthIdentifier() any   { return "u1" }
func (*probeUser) GetAuthPassword() string  { return "unused-test-hash" }
func (*probeUser) GetRememberToken() string { return "" }
func (*probeUser) SetRememberToken(string)  {}

func (u *countingUsers) FindByID(any) (auth.Authenticatable, error) { return &probeUser{}, nil }
func (u *countingUsers) FindByIDCtx(context.Context, any) (auth.Authenticatable, error) {
	return &probeUser{}, nil
}
func (u *countingUsers) FindByCredentials(map[string]any) (auth.Authenticatable, error) {
	return &probeUser{}, nil
}
func (u *countingUsers) FindByCredentialsCtx(context.Context, map[string]any) (auth.Authenticatable, error) {
	return &probeUser{}, nil
}
func (u *countingUsers) ValidateCredentials(auth.Authenticatable, map[string]any) bool {
	u.mu.Lock()
	u.checks++
	u.mu.Unlock()
	return false
}
func (u *countingUsers) UpdateRememberToken(auth.Authenticatable, string) error { return nil }
func (u *countingUsers) UpdateRememberTokenCtx(context.Context, auth.Authenticatable, string) error {
	return nil
}
func (u *countingUsers) verified() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return u.checks
}

// installCountingUsers points the default JWT scheme at the counting
// store after bootstrap (AppModule.Init installs the ORM-backed store).
func installCountingUsers(t *testing.T, a *velocity.App) *countingUsers {
	t.Helper()
	users := &countingUsers{}
	def, err := a.Auth.(*auth.Manager).DefaultScheme()
	if err != nil {
		t.Fatalf("DefaultScheme: %v", err)
	}
	scheme, ok := def.(*schemes.JWTScheme)
	if !ok {
		t.Fatalf("default scheme is %T, want *schemes.JWTScheme", def)
	}
	scheme.SetUserStore(users)
	return users
}

// loginAttempt submits wrong credentials to the probe login route through
// the global and API middleware stacks with the given forwarding headers.
func loginAttempt(t *testing.T, a *velocity.App, headers map[string]string) {
	t.Helper()
	client := velhttp.NewTestClient(t, a.Router)
	for k, v := range headers {
		client.WithHeader(k, v)
	}
	client.PostJSON("/_probe/login", map[string]any{"email": "victim@example.test", "password": "wrong"}).
		AssertUnauthorized()
}

// TestLoginBudget_SpoofedForwardedHeadersDoNotRotateBuckets is the
// review's TestForwardedHeadersBypassLoginBudgets, inverted, through the
// application's own middleware stacks and the JWT scheme's Attempt:
// rotating X-Forwarded-For and X-Forwarded-Host from an untrusted peer
// leaves every attempt in the same throttle bucket, so the pair budget
// stops the sixth candidate before verification.
func TestLoginBudget_SpoofedForwardedHeadersDoNotRotateBuckets(t *testing.T) {
	t.Setenv("AUTH_LOGIN_MAX_ATTEMPTS", "5")
	a, _ := proxyTestApp(t, nil)
	users := installCountingUsers(t, a)

	for i := 0; i < 8; i++ {
		loginAttempt(t, a, map[string]string{
			"X-Forwarded-For":   fmt.Sprintf("203.0.113.%d", i+1),
			"X-Forwarded-Host":  "attacker.example",
			"X-Forwarded-Proto": "https",
		})
	}
	if got := users.verified(); got != 5 {
		t.Fatalf("verified %d candidates, want 5 (pair budget); forwarded headers rotated the throttle bucket", got)
	}
}

// TestLoginBudget_AttackerPrefixBeforeTrustedChainDoesNotRotateBuckets
// covers the proxied deployment: the peer is a configured proxy and the
// attacker rotates an entry ahead of the real client in the chain. The
// bucket keys off the right-most untrusted hop, so the budget still holds.
func TestLoginBudget_AttackerPrefixBeforeTrustedChainDoesNotRotateBuckets(t *testing.T) {
	t.Setenv("AUTH_LOGIN_MAX_ATTEMPTS", "5")
	a, _ := proxyTestApp(t, []string{testPeer})
	users := installCountingUsers(t, a)

	for i := 0; i < 8; i++ {
		loginAttempt(t, a, map[string]string{
			"X-Forwarded-For":  fmt.Sprintf("203.0.113.%d, 198.51.100.7, %s", i+1, testPeer),
			"X-Forwarded-Host": "attacker.example",
		})
	}
	if got := users.verified(); got != 5 {
		t.Fatalf("verified %d candidates, want 5 (pair budget); the attacker prefix rotated the throttle bucket", got)
	}
}
