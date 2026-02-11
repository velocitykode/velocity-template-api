package routes

import (
	"{{MODULE_NAME}}/internal/handlers"
	"{{MODULE_NAME}}/internal/middleware"

	"github.com/velocitykode/velocity"
	"github.com/velocitykode/velocity/pkg/router"
)

func Register(v *velocity.App) {
	r := v.Router

	r.Get("/api/health", handlers.Health)
	r.Post("/api/users", handlers.CreateUser)

	r.Group("/api", func(api router.Router) {
		api.Get("/users", handlers.ListUsers)
		api.Get("/users/:id", handlers.GetUser)
		api.Get("/me", handlers.GetCurrentUser)
	}).Use(middleware.Auth)
}
