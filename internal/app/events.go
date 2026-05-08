package app

import (
	"context"

	"github.com/velocitykode/velocity/cache"
	"github.com/velocitykode/velocity/events"
	"github.com/velocitykode/velocity/log"
	"github.com/velocitykode/velocity/orm"
	"github.com/velocitykode/velocity/router"
)

// listenerFunc adapts a plain function to the events.Listener interface.
// The framework propagates a context.Context into every listener so
// downstream calls (db, cache, http) can carry deadlines and trace
// metadata; listeners that don't need it can ignore the parameter.
type listenerFunc func(ctx context.Context, event interface{}) error

func (f listenerFunc) Handle(ctx context.Context, event interface{}) error {
	return f(ctx, event)
}
func (f listenerFunc) ShouldQueue() bool { return false }

// on registers a function listener on the given dispatcher.
func on(d events.Dispatcher, event string, fn func(ctx context.Context, event interface{}) error) {
	d.Listen(event, listenerFunc(fn))
}

// Events registers event listeners for framework observability.
// main.go calls Events(v.Log) once at bootstrap and passes the returned
// callback to v.Events(...). The closure captures the logger so listeners
// can call it without reaching back through the app.
//
// Customize these listeners to add your own logging, metrics, or tracing.
func Events(logger log.Logger) func(events.Dispatcher) {
	return func(d events.Dispatcher) {
		// Request lifecycle events
		on(d, "request.started", func(_ context.Context, e interface{}) error {
			if req, ok := e.(*router.RequestStarted); ok {
				logger.Debug("Request started",
					"request_id", req.RequestID,
					"method", req.Method,
					"path", req.Path,
				)
			}
			return nil
		})

		on(d, "request.handled", func(_ context.Context, e interface{}) error {
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

		on(d, "request.failed", func(_ context.Context, e interface{}) error {
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
		on(d, "query.executed", func(_ context.Context, e interface{}) error {
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
		on(d, "cache.hit", func(_ context.Context, e interface{}) error {
			if c, ok := e.(*cache.CacheHit); ok {
				logger.Debug("Cache hit", "key", c.Key)
			}
			return nil
		})

		on(d, "cache.miss", func(_ context.Context, e interface{}) error {
			if c, ok := e.(*cache.CacheMiss); ok {
				logger.Debug("Cache miss", "key", c.Key)
			}
			return nil
		})
	}
}
