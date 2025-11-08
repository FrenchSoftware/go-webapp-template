package router

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

type Config struct {
	EnableCORS        bool          // Enable CORS middleware
	AllowedOrigins    []string      // CORS allowed origins
	EnableCompression bool          // Enable gzip compression
	EnableRateLimit   bool          // Enable rate limiting
	RateLimitRPS      int           // Requests per second (if rate limiting enabled)
	RateLimitBurst    int           // Burst size for rate limiter
	RequestTimeout    time.Duration // Request timeout duration
	LogRequests       bool          // Enable request logging
}

// DefaultConfig returns a production-ready default configuration
func DefaultConfig() Config {
	return Config{
		EnableCORS:        true,
		AllowedOrigins:    []string{"*"}, // Configure this for production!
		EnableCompression: true,
		EnableRateLimit:   true,
		RateLimitRPS:      100,
		RateLimitBurst:    200,
		RequestTimeout:    30 * time.Second,
		LogRequests:       true,
	}
}

// New creates a new gorilla/mux router with production-ready middleware
func New(cfg Config) *Router {
	r := mux.NewRouter()

	// Store config for later use
	router := WrapRouter(r)
	router.config = cfg

	// Build middleware chain using alice for clean composition
	middleware := buildMiddlewareChain(cfg)

	// Apply global middleware
	r.Use(middleware.Then)

	return router
}

// NewWithHotReload creates a router and a separate mux for hot reload endpoints
// that bypass problematic middleware (timeout, compression) for WebSocket connections
func NewWithHotReload() (*Router, *mux.Router) {
	mainRouter := mux.NewRouter()
	hotReloadRouter := mux.NewRouter()

	cfg := DefaultConfig()

	// Build middleware chain - only for main router
	middleware := buildMiddlewareChain(cfg)
	mainRouter.Use(middleware.Then)

	// Wrap the main router for route registration
	wrappedRouter := &Router{
		Router:       mainRouter,
		config:       cfg,
		hotReloadMux: hotReloadRouter,
	}

	return wrappedRouter, hotReloadRouter
}

// buildMiddlewareChain creates a composable middleware chain
func buildMiddlewareChain(cfg Config) alice.Chain {
	chain := alice.New()

	// Recovery middleware - must be first to catch panics
	chain = chain.Append(recoveryMiddleware())

	// Request ID middleware
	chain = chain.Append(requestIDMiddleware())

	// Structured logging middleware
	if cfg.LogRequests {
		chain = chain.Append(structuredLoggerMiddleware())
	}

	// Security headers middleware
	chain = chain.Append(securityHeadersMiddleware())

	// CORS middleware
	if cfg.EnableCORS {
		chain = chain.Append(corsMiddleware(cfg.AllowedOrigins))
	}

	// Compression middleware
	if cfg.EnableCompression {
		chain = chain.Append(handlers.CompressHandler)
	}

	// Rate limiting middleware
	if cfg.EnableRateLimit {
		chain = chain.Append(rateLimitMiddleware(cfg.RateLimitRPS, cfg.RateLimitBurst))
	}

	// Timeout middleware
	if cfg.RequestTimeout > 0 {
		chain = chain.Append(timeoutMiddleware(cfg.RequestTimeout))
	}

	return chain
}

// corsMiddleware creates a CORS middleware handler
func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposedHeaders:   []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           int((12 * time.Hour).Seconds()),
	})

	return c.Handler
}

// rateLimitMiddleware creates a rate limiting middleware using token bucket algorithm
func rateLimitMiddleware(rps, burst int) func(http.Handler) http.Handler {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, `{"error": "rate limit exceeded"}`, http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// timeoutMiddleware creates a timeout middleware that cancels requests exceeding the timeout
func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				if ctx.Err() == context.DeadlineExceeded {
					slog.Warn("request timeout", "path", r.URL.Path, "method", r.Method)
					http.Error(w, `{"error": "request timeout"}`, http.StatusRequestTimeout)
				}
			}
		})
	}
}

// securityHeadersMiddleware adds security headers to all responses
func securityHeadersMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' https://unpkg.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' https://*.googleusercontent.com data:; connect-src 'self' ws://localhost:* wss://localhost:* https://unpkg.com")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			next.ServeHTTP(w, r)
		})
	}
}
