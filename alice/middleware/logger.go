package middleware

import (
	"context"
	"net/http"

	"github.com/promoboxx/go-service/alice/middleware/lrw"
	"github.com/sirupsen/logrus"
)

const (
	logFieldRequestID  = "request_id"
	logFieldUserID     = "user_id"
	logFieldMethod     = "method"
	logFieldStatusCode = "status_code"
	logFieldPath       = "path"
	logFieldError      = "error"
)

// Logger injects a logger into the context
type Logger interface {
	Log(h http.Handler) http.Handler
}

type logger struct {
	entry       *logrus.Entry
	logRequests bool
}

// NewLogrusLogger allows you to setup a base entry to use for logging
func NewLogrusLogger(baseEntry *logrus.Entry, logRequests bool) Logger {
	return &logger{entry: baseEntry, logRequests: logRequests}
}

func (l *logger) Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingResponseWriter := &lrw.LoggingResponseWriter{w, http.StatusOK, nil}

		fields := logrus.Fields{
			logFieldRequestID: getRequestIDFromContext(r.Context()),
			logFieldUserID:    getInsecureUserIDFromContext(r.Context()),
			logFieldMethod:    r.Method,
			logFieldPath:      r.URL.String(),
		}

		entry := l.entry.WithFields(fields)

		// add logger to the context
		ctx := context.WithValue(r.Context(), contextKeyLogger, entry)
		r = r.WithContext(ctx)

		h.ServeHTTP(loggingResponseWriter, r)

		fields[logFieldStatusCode] = loggingResponseWriter.StatusCode

		if loggingResponseWriter.InnerError != nil {
			fields[logFieldError] = loggingResponseWriter.InnerError
		}

		if l.logRequests {
			entry.Printf("Finished request")
		}
	})
}
