package router

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
)

// requestIDMiddleware adds a unique request ID to each request
func requestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			w.Header().Set("X-Request-ID", requestID)
			r.Header.Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// structuredLoggerMiddleware logs requests using slog (structured logging)
func structuredLoggerMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap ResponseWriter to capture status code
			lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(lrw, r)

			duration := time.Since(start)

			slog.Info("http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", lrw.statusCode,
				"duration", duration,
				"ip", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"request_id", r.Header.Get("X-Request-ID"),
			)
		})
	}
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// recoveryMiddleware recovers from panics and logs them
func recoveryMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					slog.Error("panic recovered",
						"error", err,
						"path", r.URL.Path,
						"method", r.Method,
						"request_id", r.Header.Get("X-Request-ID"),
						"stack", string(debug.Stack()),
					)

					http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// NoCache is a middleware that adds headers to prevent client caching
func NoCache() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			next.ServeHTTP(w, r)
		})
	}
}

// APIVersion is a middleware that adds an API version header to responses
func APIVersion(version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-API-Version", version)
			next.ServeHTTP(w, r)
		})
	}
}
