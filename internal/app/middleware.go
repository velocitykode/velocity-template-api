package app

import (
	"{{MODULE_NAME}}/internal/middleware"

	"github.com/velocitykode/velocity"
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
		middleware.LoggingMiddleware,                          // Log all requests
		middleware.TrustProxiesMiddleware,                     // Handle X-Forwarded-* headers
		middleware.CORSMiddleware,                             // CORS preflight + headers
		middleware.PreventRequestsDuringMaintenanceMiddleware, // 503 when in maintenance mode
		middleware.ValidatePostSizeMiddleware(10<<20),         // Reject requests > 10MB
		middleware.TrimStringsMiddleware,                      // Trim whitespace from string inputs
		middleware.ConvertEmptyStringsToNullMiddleware,        // Convert "" to nil
	)

	m.API(
		middleware.EnsureJSONMiddleware, // Force JSON content-type on responses
	)
}
