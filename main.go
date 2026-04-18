package main

import (
	"log"

	"{{MODULE_NAME}}/internal/app"
	"{{MODULE_NAME}}/internal/commands"
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

	if err := v.
		Providers(app.Configure).
		Middleware(app.Middleware).
		Routes(routes.Register).
		Commands(commands.Register).
		Events(app.Events(v.Log)).
		Serve(); err != nil {
		log.Fatal(err)
	}
}
