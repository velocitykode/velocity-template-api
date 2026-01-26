package app

import (
	"{{MODULE_NAME}}/internal/middleware"

	"github.com/velocitykode/velocity/pkg/router"
)

// MiddlewareStacks defines all middleware stacks for the application
type MiddlewareStacks struct {
	Global []router.MiddlewareFunc
	API    []router.MiddlewareFunc
}

// GetMiddlewareStacks returns configured middleware stacks
func GetMiddlewareStacks() *MiddlewareStacks {
	return &MiddlewareStacks{
		Global: globalMiddleware(),
		API:    apiMiddleware(),
	}
}

// globalMiddleware returns middleware that runs for ALL requests
func globalMiddleware() []router.MiddlewareFunc {
	return []router.MiddlewareFunc{
		middleware.RecoveryMiddleware,                         // Catch panics and return 500
		middleware.LoggingMiddleware,                          // Log all requests
		middleware.TrustProxiesMiddleware,                     // Handle X-Forwarded-* headers
		middleware.CORSMiddleware,                             // Handle CORS preflight and headers
		middleware.PreventRequestsDuringMaintenanceMiddleware, // Return 503 when in maintenance mode
		middleware.ValidatePostSizeMiddleware(10 << 20),       // Reject requests > 10MB
		middleware.TrimStringsMiddleware,                      // Trim whitespace from string inputs
		middleware.ConvertEmptyStringsToNullMiddleware,        // Convert "" to nil for cleaner handling
	}
}

// apiMiddleware returns middleware for API requests (JSON APIs)
func apiMiddleware() []router.MiddlewareFunc {
	return []router.MiddlewareFunc{
		middleware.EnsureJSONMiddleware, // Ensure response is JSON formatted
	}
}
