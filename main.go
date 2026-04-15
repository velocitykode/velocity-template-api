package main

import (
	"log"
	"os"

	"{{MODULE_NAME}}/internal/app"
	"{{MODULE_NAME}}/routes"

	"github.com/velocitykode/velocity"

	// Blank import so each migration file's init() runs and calls
	// migrate.Register() - otherwise `vel migrate` finds nothing.
	_ "{{MODULE_NAME}}/database/migrations"
)

func main() {
	v, err := velocity.New()
	if err != nil {
		log.Fatal(err)
	}

	chain := v.
		Providers(app.Configure).
		Middleware(app.Middleware).
		Routes(routes.Register).
		Events(app.Events(v.Log))

	// With CLI args (`vel migrate`, `vel make:handler`, ...) dispatch
	// the command. Routes are still registered above so `vel route:list`
	// sees them.
	if len(os.Args) > 1 {
		if err := chain.Run(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// No args - start the HTTP server.
	if err := chain.Serve(); err != nil {
		log.Fatal(err)
	}
}
