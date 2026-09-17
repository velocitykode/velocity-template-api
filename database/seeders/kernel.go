package seeders

import (
	"github.com/velocitykode/velocity"
)

// Register declares the application's database seeders in the order they
// run: list a seeder after the ones whose rows it depends on. Generated
// seeders (via `vel gen seeder`) are added here; `vel db seed` runs the list
// and `vel db seed --only <name>` runs a single entry.
func Register(r *velocity.Seeders) {
	// r.Add(&RoleSeeder{})
}
