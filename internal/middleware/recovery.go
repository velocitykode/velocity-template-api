package middleware

import (
	"net/http"

	"github.com/velocitykode/velocity/pkg/log"
	"github.com/velocitykode/velocity/pkg/router"
)

func RecoveryMiddleware(next router.HandlerFunc) router.HandlerFunc {
	return func(c *router.Context) error {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered", "error", err)
				http.Error(c.Response, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		return next(c)
	}
}
