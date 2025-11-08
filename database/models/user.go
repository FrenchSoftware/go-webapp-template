package models

import "time"

type User struct {
	ID            int64     `json:"id"`
	GoogleID      string    `json:"google_id"`
	Email         string    `json:"email"`
	Name          string    `json:"name"`
	GivenName     *string   `json:"given_name,omitempty"`
	FamilyName    *string   `json:"family_name,omitempty"`
	Picture       *string   `json:"picture,omitempty"`
	Locale        *string   `json:"locale,omitempty"`
	VerifiedEmail bool      `json:"verified_email"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Session struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
