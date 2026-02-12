package routes

import (
	"{{MODULE_NAME}}/internal/handlers"

	"github.com/velocitykode/velocity"
)

func Register(v *velocity.App) {
	r := v.Router

	r.Get("/api/health", handlers.Health)
}
