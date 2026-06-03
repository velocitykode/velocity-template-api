package app

import (
	"os"
	"strings"

	"{{MODULE_NAME}}/internal/middleware"

	"github.com/velocitykode/velocity"
	"github.com/velocitykode/velocity/router"
)

// Middleware configures the application's middleware stacks.
//
// The framework calls this once during bootstrap with a *MiddlewareStack
// that splits middleware into two scopes for an API-only project:
//
//   - Global: runs on every request
//   - API:    runs on routes inside r.API(prefix, ...)
//
// (Web is also available on the stack - unused here because the API
// template doesn't ship browser-rendered routes.)
func Middleware(m *velocity.MiddlewareStack) {
	m.Global(
		middleware.LoggingMiddleware,      // Log all requests (no framework export yet)
		middleware.TrustProxiesMiddleware, // Handle X-Forwarded-* headers (no framework export yet)
		router.CORS(router.CORSConfig{ // Framework CORS (velocity/router/cors.go)
			AllowedOrigins:   corsOrigins(),
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
			AllowCredentials: true,
		}),
		velocity.PreventRequestsDuringMaintenance(),    // Framework maintenance gate (velocity/maintenance.go)
		router.BodyLimit(10<<20),                       // Framework body-limit (velocity/router/body_limit.go) - 10MB
		middleware.TrimStringsMiddleware,               // Trim whitespace from string inputs
		middleware.ConvertEmptyStringsToNullMiddleware, // Convert "" to nil
	)

	m.API(
		middleware.EnsureJSONMiddleware, // Force JSON response content-type (sets response header; not the same as router.ContentTypeJSON which validates request headers)
	)
}

// corsOrigins returns the CORS allowlist from the CORS_ALLOWED_ORIGINS env var
// (comma-separated). It falls back to APP_URL when unset, and never returns
// "*": a wildcard combined with AllowCredentials echoes the request origin
// back, letting any site make credentialed requests. Set CORS_ALLOWED_ORIGINS
// to your frontend origin(s) for cross-origin deployments.
func corsOrigins() []string {
	if raw := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); raw != "" {
		var origins []string
		for _, o := range strings.Split(raw, ",") {
			if o = strings.TrimSpace(o); o != "" {
				origins = append(origins, o)
			}
		}
		if len(origins) > 0 {
			return origins
		}
	}
	if appURL := strings.TrimSpace(os.Getenv("APP_URL")); appURL != "" {
		return []string{appURL}
	}
	// Nothing configured: return no origins so the framework rejects all
	// cross-origin requests (fail closed). Same-origin requests are unaffected.
	return nil
}
