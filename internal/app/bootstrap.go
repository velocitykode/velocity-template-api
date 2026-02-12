package app

import (
	"{{MODULE_NAME}}/config"

	"github.com/velocitykode/velocity"
	"github.com/velocitykode/velocity/pkg/auth"
	"github.com/velocitykode/velocity/pkg/auth/drivers/guards"
)

// Bootstrap configures app-specific services on the Velocity app instance.
// Core services (crypto, ORM, logger, cache, events, queue, storage) are
// already initialized by velocity.Default().
func Bootstrap(v *velocity.App) error {
	// 1. Register auth guards (app-specific: session guard with user model)
	if err := bootstrapAuth(v); err != nil {
		return err
	}

	// 2. Register event listeners (app-specific)
	initEvents(v)

	// 3. Apply middleware to the router
	bootstrapMiddleware(v)

	return nil
}

func bootstrapAuth(v *velocity.App) error {
	authManager := v.Auth.(*auth.Manager)
	sessionConfig := auth.NewSessionConfigFromEnv()
	provider := auth.NewORMUserProvider(v.DB.DB(), config.GetAuthModel(), authManager.GetHasher())
	sessionGuard, err := guards.NewSessionGuard(provider, sessionConfig, v.Crypto)
	if err != nil {
		return err
	}

	authManager.RegisterGuard(config.GetAuthGuard(), sessionGuard)
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
