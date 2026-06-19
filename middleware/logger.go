package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

type contextKey string

const errorKey contextKey = "handler_error"

type ErrorContainer struct {
	Err error
}

func AddError(r *http.Request, err error) {
	if container, ok := r.Context().Value(errorKey).(*ErrorContainer); ok {
		container.Err = err
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		errContainer := &ErrorContainer{}
		r = r.WithContext(context.WithValue(r.Context(), errorKey, errContainer))

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		log.Printf("[HTTP] Started %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		if rw.statusCode >= 400 {
			errStr := "Unknown error"
			if errContainer.Err != nil {
				errStr = errContainer.Err.Error()
			}
			log.Printf("[HTTP] Completed %s %s | Status: %d (ERROR) | Error: %s | Duration: %s",
				r.Method, r.URL.Path, rw.statusCode, errStr, duration)
		} else {
			log.Printf("[HTTP] Completed %s %s | Status: %d (SUCCESS) | Duration: %s",
				r.Method, r.URL.Path, rw.statusCode, duration)
		}
	})
}
