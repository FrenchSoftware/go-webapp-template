package controllers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hyperstitieux/template/auth"
	"github.com/hyperstitieux/template/database/models"
	"github.com/hyperstitieux/template/database/repositories"
	"golang.org/x/oauth2"
)

const (
	stateCookieName    = "oauth_state"
	redirectCookieName = "oauth_redirect"
	sessionDuration    = 30 * 24 * time.Hour
)

type GoogleOAuthController interface {
	Redirect(w http.ResponseWriter, r *http.Request) error
	Callback(w http.ResponseWriter, r *http.Request) error
}

type googleOAuthController struct {
	users       repositories.UsersRepository
	oauthConfig *oauth2.Config
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func NewGoogleOAuthController(users repositories.UsersRepository, oauthConfig *oauth2.Config) GoogleOAuthController {
	return &googleOAuthController{
		users:       users,
		oauthConfig: oauthConfig,
	}
}

// Redirect initiates the OAuth2 flow by redirecting to Google
func (c *googleOAuthController) Redirect(w http.ResponseWriter, r *http.Request) error {
	// Generate a random state token
	state, err := generateRandomToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate state token: %w", err)
	}

	// Store state in a cookie for verification in callback
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	// Store redirect URL from query parameter (where to go after auth)
	redirectTo := r.URL.Query().Get("redirect")
	if redirectTo == "" {
		redirectTo = "/" // Default to home page
	}
	http.SetCookie(w, &http.Cookie{
		Name:     redirectCookieName,
		Value:    redirectTo,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to Google's OAuth2 consent page
	url := c.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

// Callback handles the OAuth2 callback from Google
func (c *googleOAuthController) Callback(w http.ResponseWriter, r *http.Request) error {
	// Verify state token
	stateCookie, err := r.Cookie(stateCookieName)
	if err != nil {
		return fmt.Errorf("state cookie not found: %w", err)
	}

	state := r.URL.Query().Get("state")
	if state != stateCookie.Value {
		return fmt.Errorf("invalid state token")
	}

	// Get redirect URL from cookie
	redirectTo := "/"
	if redirectCookie, err := r.Cookie(redirectCookieName); err == nil {
		redirectTo = redirectCookie.Value
	}

	// Clear state and redirect cookies
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     redirectCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		return fmt.Errorf("authorization code not found")
	}

	// Exchange code for token
	token, err := c.oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from Google
	userInfo, err := c.getUserInfo(token)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if user exists
	user, err := c.users.GetUserByGoogleID(userInfo.ID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Create or update user
	if user == nil {
		// Create new user
		user = &models.User{
			GoogleID:      userInfo.ID,
			Email:         userInfo.Email,
			Name:          userInfo.Name,
			GivenName:     stringPtr(userInfo.GivenName),
			FamilyName:    stringPtr(userInfo.FamilyName),
			Picture:       stringPtr(userInfo.Picture),
			Locale:        stringPtr(userInfo.Locale),
			VerifiedEmail: userInfo.VerifiedEmail,
		}
		if err := c.users.CreateUser(user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Update existing user
		user.Email = userInfo.Email
		user.Name = userInfo.Name
		user.GivenName = stringPtr(userInfo.GivenName)
		user.FamilyName = stringPtr(userInfo.FamilyName)
		user.Picture = stringPtr(userInfo.Picture)
		user.Locale = stringPtr(userInfo.Locale)
		user.VerifiedEmail = userInfo.VerifiedEmail
		if err := c.users.UpdateUser(user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Create session
	sessionToken, err := generateRandomToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate session token: %w", err)
	}

	session := &models.Session{
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(sessionDuration),
	}
	if err := c.users.CreateSession(session); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Set session cookie using secure cookie helper
	auth.SetSessionCookie(w, r, sessionToken, sessionDuration)

	// Redirect to original page or home page
	http.Redirect(w, r, redirectTo, http.StatusTemporaryRedirect)
	return nil
}

// getUserInfo fetches user information from Google using the access token
func (c *googleOAuthController) getUserInfo(token *oauth2.Token) (*GoogleUserInfo, error) {
	client := c.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: %s", string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}

// generateRandomToken generates a secure random token
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// stringPtr returns a pointer to a string
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
