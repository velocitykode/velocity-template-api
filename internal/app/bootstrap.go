package app

import (
	"github.com/velocitykode/velocity"
)

// Bootstrap configures app-specific services on the Velocity app instance.
func Bootstrap(v *velocity.App) error {
	// 1. Register event listeners
	initEvents(v)

	// 2. Apply middleware to the router
	bootstrapMiddleware(v)

	return nil
}

func bootstrapMiddleware(v *velocity.App) {
	stacks := GetMiddlewareStacks()

	for _, mw := range stacks.Global {
		v.Router.Use(mw)
	}
	for _, mw := range stacks.API {
		v.Router.Use(mw)
	}
}
