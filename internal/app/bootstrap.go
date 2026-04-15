package app

import (
	"github.com/velocitykode/velocity"
)

// Configure registers the app's service providers. main.go passes this
// to v.Providers(...). Add your own providers here as the app grows
// (e.g. JWT auth setup, third-party SDK wiring).
//
// Empty by default — the API template ships with no app-specific
// services to register. The framework's own services (DB, cache,
// queue, logger, etc.) are initialized by velocity.New().
func Configure(_ *velocity.ProviderRegistry) {
	// reg.Add(&AuthProvider{}, &MyService{})
}
