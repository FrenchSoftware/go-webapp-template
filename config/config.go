package config

import (
	"github.com/hyperstitieux/template/env"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type config struct {
	HTTPAddr         string
	DatabaseURL      string
	GoogleOAuthConfig *oauth2.Config
	BaseURL          string
}

type Config *config

func New() Config {
	baseURL := env.GetVar("BASE_URL", "http://localhost:8080")

	return &config{
		HTTPAddr:    env.GetVar("HTTP_ADDR", ":8080"),
		DatabaseURL: env.GetVar("DATABASE_URL", "file:app.db"),
		BaseURL:     baseURL,
		GoogleOAuthConfig: &oauth2.Config{
			ClientID:     env.GetVar("GOOGLE_CLIENT_ID", ""),
			ClientSecret: env.GetVar("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  baseURL + "/auth/google/callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}
