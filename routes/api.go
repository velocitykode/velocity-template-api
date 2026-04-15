package routes

import (
	"{{MODULE_NAME}}/internal/handlers"

	"github.com/velocitykode/velocity"
	"github.com/velocitykode/velocity/router"
)

// Register defines all application routes. main.go passes this function
// to v.Routes(...). The framework calls it with a *velocity.Routing
// already wired with the configured middleware stacks.
func Register(r *velocity.Routing) {
	// Operational endpoint sits at the top level — no middleware stack
	// so load balancers can probe /health cheaply.
	r.Health("/health")

	// API routes — every route here is prefixed with /api and runs the
	// API middleware stack (JSON enforcement, etc.).
	r.API("/api", func(api router.Router) {
		api.Get("/health", handlers.Health) // /api/health
	})
}
