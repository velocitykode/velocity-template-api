package app

import (
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
			AllowedOrigins:   []string{"*"},
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
