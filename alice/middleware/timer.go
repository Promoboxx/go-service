package middleware

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/opentracing/opentracing-go"
)

// Timer can time a handler and log it
type Timer interface {
	Time(name string) alice.Constructor
}

type nullTimer struct{}

func (n *nullTimer) Time(name string) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}
}

func NewNullTimer() Timer {
	return &nullTimer{}
}

type openTracingTimer struct {
}

// NewOpenTracingTimer creates a new timer that uses opentracing spans
func NewOpenTracingTimer() Timer {
	return &openTracingTimer{}
}

// Time returns a middleware that will handle creating spans and tracing
// for each request
func (ott *openTracingTimer) Time(name string) alice.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			span, ctx := opentracing.StartSpanFromContext(ctx, name)
			defer span.Finish()
		})
	}
}
