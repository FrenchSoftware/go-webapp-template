package auth

import (
	"log/slog"
	"net/http"

	"github.com/hyperstitieux/template/database/repositories"
)

const (
	// SessionCookieName is the name of the session cookie
	SessionCookieName = "session"
)

// AuthMiddleware creates a middleware that authenticates requests using session cookies
func AuthMiddleware(users repositories.UsersRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try to get session cookie
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				slog.Debug("no session cookie found",
					"path", r.URL.Path,
					"error", err,
				)
				// No session cookie - continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			slog.Debug("session cookie found",
				"path", r.URL.Path,
				"cookie_name", SessionCookieName,
			)

			// Validate session and get user
			user, err := users.GetUserBySessionToken(cookie.Value)
			if err != nil {
				slog.Error("failed to get user by session token",
					"error", err,
					"path", r.URL.Path,
				)
				// Invalid session - continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			// Session not found or expired
			if user == nil {
				slog.Debug("session not found or expired",
					"path", r.URL.Path,
				)
				// Continue without authentication
				next.ServeHTTP(w, r)
				return
			}

			slog.Debug("user authenticated",
				"path", r.URL.Path,
				"user_id", user.ID,
				"user_email", user.Email,
			)

			// Attach user to request context
			r = SetCurrentUser(r, user)

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuthMiddleware creates a middleware that requires authentication
// If the user is not authenticated, it redirects to the login page
func RequireAuthMiddleware(redirectURL string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !IsAuthenticated(r) {
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAuthHandler wraps a handler function that requires authentication
func RequireAuthHandler(redirectURL string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !IsAuthenticated(r) {
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}

		handler(w, r)
	}
}
