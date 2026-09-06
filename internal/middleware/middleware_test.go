package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/velocitykode/velocity/router"
	"github.com/velocitykode/velocity/velocitytest"
)

// newTestCtx constructs a fully-wired router.Context backed by a
// velocitytest.NewApp() so that ctx.Log() and any other service-touching
// methods don't panic. The returned recorder is the same one assigned as
// ctx.Response, so test assertions can read headers/status from it.
func newTestCtx(t *testing.T, method, path string) (*router.Context, *httptest.ResponseRecorder) {
	t.Helper()
	app, err := velocitytest.NewApp()
	if err != nil {
		t.Fatalf("velocitytest.NewApp: %v", err)
	}
	t.Cleanup(func() { _ = app.Shutdown(context.Background()) })

	ctx, rec := router.NewTestContext(method, path)
	ctx.SetServices(app.Services)
	return ctx, rec
}

func TestEnsureJSONMiddleware(t *testing.T) {
	handler := EnsureJSONMiddleware(func(c *router.Context) error {
		return nil
	})

	ctx, rec := newTestCtx(t, http.MethodGet, "/")

	if err := handler(ctx); err != nil {
		t.Fatalf("EnsureJSONMiddleware() returned error: %v", err)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("EnsureJSONMiddleware() Content-Type = %q, want %q", contentType, "application/json")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	nextCalled := false
	handler := LoggingMiddleware(func(c *router.Context) error {
		nextCalled = true
		return nil
	})

	ctx, _ := newTestCtx(t, http.MethodGet, "/test")

	if err := handler(ctx); err != nil {
		t.Fatalf("LoggingMiddleware() returned error: %v", err)
	}

	if !nextCalled {
		t.Error("LoggingMiddleware() did not call next handler")
	}
}
