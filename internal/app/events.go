package app

import (
	"github.com/velocitykode/velocity"
	"github.com/velocitykode/velocity/pkg/cache"
	"github.com/velocitykode/velocity/pkg/events"
	"github.com/velocitykode/velocity/pkg/orm"
	"github.com/velocitykode/velocity/pkg/router"
)

// listenerFunc adapts a plain function to the events.Listener interface.
type listenerFunc func(event interface{}) error

func (f listenerFunc) Handle(event interface{}) error { return f(event) }
func (f listenerFunc) ShouldQueue() bool              { return false }

// on registers a function listener on the given dispatcher.
func on(d events.Dispatcher, event string, fn func(event interface{}) error) {
	d.Listen(event, listenerFunc(fn))
}

// initEvents registers event listeners for framework observability.
// Customize these listeners to add your own logging, metrics, or tracing.
func initEvents(v *velocity.App) {
	logger := v.Log

	// Request lifecycle events
	on(v.Events, "request.started", func(e interface{}) error {
		if req, ok := e.(*router.RequestStarted); ok {
			logger.Debug("Request started",
				"request_id", req.RequestID,
				"method", req.Method,
				"path", req.Path,
			)
		}
		return nil
	})

	on(v.Events, "request.handled", func(e interface{}) error {
		if req, ok := e.(*router.RequestHandled); ok {
			logger.Info("Request completed",
				"request_id", req.RequestID,
				"method", req.Method,
				"path", req.Path,
				"status", req.StatusCode,
				"duration", req.Duration,
			)
		}
		return nil
	})

	on(v.Events, "request.failed", func(e interface{}) error {
		if req, ok := e.(*router.RequestFailed); ok {
			logger.Error("Request failed",
				"request_id", req.RequestID,
				"error", req.Error,
				"recovered", req.Recovered,
			)
		}
		return nil
	})

	// Database query events
	on(v.Events, "query.executed", func(e interface{}) error {
		if q, ok := e.(*orm.QueryExecuted); ok {
			logger.Debug("Query executed",
				"sql", q.SQL,
				"duration", q.Duration,
				"rows", q.RowsAffected,
			)
		}
		return nil
	})

	// Cache events
	on(v.Events, "cache.hit", func(e interface{}) error {
		if c, ok := e.(*cache.CacheHit); ok {
			logger.Debug("Cache hit", "key", c.Key)
		}
		return nil
	})

	on(v.Events, "cache.miss", func(e interface{}) error {
		if c, ok := e.(*cache.CacheMiss); ok {
			logger.Debug("Cache miss", "key", c.Key)
		}
		return nil
	})
}
