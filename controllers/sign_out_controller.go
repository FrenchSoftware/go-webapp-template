package controllers

import (
	"fmt"
	"net/http"

	"github.com/hyperstitieux/template/auth"
	"github.com/hyperstitieux/template/database/repositories"
)

type SignOutController interface {
	Handle(w http.ResponseWriter, r *http.Request) error
}

type signOutController struct {
	users repositories.UsersRepository
}

func NewSignOutController(users repositories.UsersRepository) SignOutController {
	return &signOutController{
		users: users,
	}
}

// Handle handles user sign out by deleting the session and clearing the cookie
func (c *signOutController) Handle(w http.ResponseWriter, r *http.Request) error {
	// Get session token from cookie
	token, err := auth.GetSessionToken(r)
	if err != nil {
		// No session cookie - just redirect
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil
	}

	// Delete session from database
	if err := c.users.DeleteSession(token); err != nil {
		// Log error but continue with logout
		fmt.Printf("failed to delete session: %v\n", err)
	}

	// Clear session cookie
	auth.ClearSessionCookie(w, r)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return nil
}
