package routes

import (
	"{{MODULE_NAME}}/internal/handlers"
	"{{MODULE_NAME}}/internal/middleware"

	"github.com/velocitykode/velocity/pkg/router"
)

func init() {
	router.Register(func(r router.Router) {
		// Health check (public)
		r.Get("/api/health", handlers.Health)

		// Public routes
		r.Post("/api/users", handlers.CreateUser)

		// Protected routes (require authentication)
		r.Group("/api", func(api router.Router) {
			api.Get("/users", handlers.ListUsers)
			api.Get("/users/:id", handlers.GetUser)
			api.Get("/me", handlers.GetCurrentUser)
		}).Use(middleware.Auth)
	})
}
