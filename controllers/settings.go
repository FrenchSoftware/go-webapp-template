package controllers

import (
	"log/slog"
	"net/http"

	"github.com/frenchsoftware/libvalidator/validator"
	"github.com/hyperstitieux/template/auth"
	"github.com/hyperstitieux/template/database/repositories"
	"github.com/hyperstitieux/template/pages"
)

type SettingsController struct {
	users repositories.UsersRepository
}

func NewSettingsController(users repositories.UsersRepository) *SettingsController {
	return &SettingsController{
		users: users,
	}
}

// UpdateProfile handles profile update requests
func (c *SettingsController) UpdateProfile(w http.ResponseWriter, r *http.Request) error {
	// Get authenticated user
	user := auth.GetCurrentUser(r)
	if user == nil {
		http.Redirect(w, r, "/auth/google?redirect=/settings", http.StatusTemporaryRedirect)
		return nil
	}

	// Validate form data
	v := validator.New(
		validator.Field("name").Required().MinLength(1).MaxLength(50),
	)

	ok, errs := v.Validate(r)
	if !ok {
		// Render settings page with validation errors
		return pages.SettingsWithErrors(w, r, errs)
	}

	// Get name from form
	name := r.FormValue("name")

	// Update user name
	user.Name = name
	if err := c.users.UpdateUser(user); err != nil {
		slog.Error("failed to update user", "error", err, "user_id", user.ID)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return nil
	}

	// Redirect back to settings page
	http.Redirect(w, r, "/settings", http.StatusSeeOther)
	return nil
}

// DeleteAccount handles account deletion requests
func (c *SettingsController) DeleteAccount(w http.ResponseWriter, r *http.Request) error {
	// Get authenticated user
	user := auth.GetCurrentUser(r)
	if user == nil {
		http.Redirect(w, r, "/auth/google", http.StatusTemporaryRedirect)
		return nil
	}

	// Delete user account
	if err := c.users.DeleteUser(user.ID); err != nil {
		slog.Error("failed to delete user", "error", err, "user_id", user.ID)
		http.Error(w, "Failed to delete account", http.StatusInternalServerError)
		return nil
	}

	// Clear session cookie
	auth.ClearSessionCookie(w, r)

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}
