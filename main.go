package main

import (
	"log"

	"{{MODULE_NAME}}/internal/app"
	"{{MODULE_NAME}}/internal/commands"
	"{{MODULE_NAME}}/routes"

	"github.com/velocitykode/velocity"

	// Standard driver bundles. Since the framework decouple, heavy drivers
	// live in opt-in leaves and no longer auto-register; blank-importing each
	// subsystem's `standard` bundle registers every driver (memory/file + redis
	// cache, modernc/postgres/mysql/sqlite orm, memory/redis/database queue,
	// local/s3 storage) so config can select any of them at runtime.
	_ "github.com/velocitykode/velocity/cache/standard"
	_ "github.com/velocitykode/velocity/orm/standard"
	_ "github.com/velocitykode/velocity/queue/standard"
	_ "github.com/velocitykode/velocity/storage/standard"

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
