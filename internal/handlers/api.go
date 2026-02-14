package handlers

import (
	"github.com/velocitykode/velocity/router"
)

// Health returns the API health status
func Health(ctx *router.Context) error {
	return ctx.JSON(router.StatusOK, map[string]string{
		"status": "healthy",
	})
}
