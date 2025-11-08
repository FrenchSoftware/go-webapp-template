package auth

import (
	"context"
	"net/http"

	"github.com/hyperstitieux/template/database/models"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UserContextKey is the key used to store the user in the request context
	UserContextKey contextKey = "user"
)

// GetCurrentUser retrieves the authenticated user from the request context
func GetCurrentUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}

// SetCurrentUser stores the user in the request context
func SetCurrentUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

// IsAuthenticated checks if the current request has an authenticated user
func IsAuthenticated(r *http.Request) bool {
	return GetCurrentUser(r) != nil
}

// RequireAuth returns true if the user is authenticated, false otherwise
func RequireAuth(r *http.Request) (*models.User, bool) {
	user := GetCurrentUser(r)
	if user == nil {
		return nil, false
	}
	return user, true
}
