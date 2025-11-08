package views

import (
	"net/http"

	"github.com/hyperstitieux/template/auth"
	"github.com/hyperstitieux/template/database/models"
)

// GetUser retrieves the authenticated user from the request context
// Returns nil if no user is authenticated
// This delegates to the auth package to ensure consistent context key usage
func GetUser(r *http.Request) *models.User {
	return auth.GetCurrentUser(r)
}
