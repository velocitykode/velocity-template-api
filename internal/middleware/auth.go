package middleware

import (
	"net/http"

	"github.com/velocitykode/velocity/pkg/auth"
	"github.com/velocitykode/velocity/pkg/router"
)

// Auth returns 401 Unauthorized if not authenticated
func Auth(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx *router.Context) error {
		if !auth.Check(ctx.Request) {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Unauthorized",
			})
		}
		return next(ctx)
	}
}
