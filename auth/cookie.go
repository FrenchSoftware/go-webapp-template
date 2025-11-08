package auth

import (
	"net/http"
	"time"
)

// CookieConfig holds configuration for secure cookie creation
type CookieConfig struct {
	Name     string
	Value    string
	Path     string
	MaxAge   int           // in seconds
	Duration time.Duration // alternative to MaxAge
	Secure   bool          // HTTPS only
	HttpOnly bool
	SameSite http.SameSite
	Domain   string
}

// DefaultCookieConfig returns a secure cookie configuration
func DefaultCookieConfig(name, value string, duration time.Duration) CookieConfig {
	return CookieConfig{
		Name:     name,
		Value:    value,
		Path:     "/",
		Duration: duration,
		Secure:   true,  // Always true in production
		HttpOnly: true,  // Prevent XSS attacks
		SameSite: http.SameSiteLaxMode, // CSRF protection
	}
}

// SetSecureCookie creates and sets a secure cookie on the response
func SetSecureCookie(w http.ResponseWriter, r *http.Request, config CookieConfig) {
	// Determine if we should use Secure flag based on environment
	secure := config.Secure

	// In development (no TLS), we may want to disable Secure flag
	// but keep it enabled in production
	// Check if running on localhost (with or without port)
	host := r.Host
	if r.TLS == nil && (host == "localhost" || host == "127.0.0.1" ||
		len(host) > 10 && host[:10] == "localhost:" ||
		len(host) > 10 && host[:10] == "127.0.0.1:") {
		secure = false
	}

	maxAge := config.MaxAge
	if maxAge == 0 && config.Duration > 0 {
		maxAge = int(config.Duration.Seconds())
	}

	cookie := &http.Cookie{
		Name:     config.Name,
		Value:    config.Value,
		Path:     config.Path,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: config.HttpOnly,
		SameSite: config.SameSite,
		Domain:   config.Domain,
	}

	http.SetCookie(w, cookie)
}

// SetSessionCookie creates and sets a session cookie
func SetSessionCookie(w http.ResponseWriter, r *http.Request, token string, duration time.Duration) {
	config := DefaultCookieConfig(SessionCookieName, token, duration)
	SetSecureCookie(w, r, config)
}

// ClearSessionCookie removes the session cookie
func ClearSessionCookie(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   r.TLS != nil,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

// GetSessionToken retrieves the session token from the request cookie
func GetSessionToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
