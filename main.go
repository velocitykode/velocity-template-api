package main

import (
	"log"

	"{{MODULE_NAME}}/internal/app"
	"{{MODULE_NAME}}/routes"

	"github.com/velocitykode/velocity"

	// Blank import so each migration file's init() runs and calls
	// migrate.Register() — otherwise `vel migrate` finds nothing.
	_ "{{MODULE_NAME}}/database/migrations"
)

func main() {
	v, err := velocity.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Bootstrap(v); err != nil {
		log.Fatal(err)
	}

	routes.Register(v)

	if err := v.Run(); err != nil {
		log.Fatal(err)
	}
}
