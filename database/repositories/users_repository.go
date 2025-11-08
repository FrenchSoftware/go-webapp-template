package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/hyperstitieux/template/database/models"
)

type UsersRepository interface {
	// User operations
	CreateUser(user *models.User) error
	GetUserByID(id int64) (*models.User, error)
	GetUserByGoogleID(googleID string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error

	// Session operations
	CreateSession(session *models.Session) error
	GetSessionByToken(token string) (*models.Session, error)
	GetUserBySessionToken(token string) (*models.User, error)
	DeleteSession(token string) error
	DeleteExpiredSessions() error
}

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) UsersRepository {
	return &usersRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *usersRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (google_id, email, name, given_name, family_name, picture, locale, verified_email)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		user.GoogleID,
		user.Email,
		user.Name,
		user.GivenName,
		user.FamilyName,
		user.Picture,
		user.Locale,
		user.VerifiedEmail,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	user.ID = id
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return nil
}

// GetUserByID retrieves a user by their ID
func (r *usersRepository) GetUserByID(id int64) (*models.User, error) {
	query := `
		SELECT id, google_id, email, name, given_name, family_name, picture, locale, verified_email, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.Name,
		&user.GivenName,
		&user.FamilyName,
		&user.Picture,
		&user.Locale,
		&user.VerifiedEmail,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

// GetUserByGoogleID retrieves a user by their Google ID
func (r *usersRepository) GetUserByGoogleID(googleID string) (*models.User, error) {
	query := `
		SELECT id, google_id, email, name, given_name, family_name, picture, locale, verified_email, created_at, updated_at
		FROM users
		WHERE google_id = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, googleID).Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.Name,
		&user.GivenName,
		&user.FamilyName,
		&user.Picture,
		&user.Locale,
		&user.VerifiedEmail,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by google id: %w", err)
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email address
func (r *usersRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, google_id, email, name, given_name, family_name, picture, locale, verified_email, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.Name,
		&user.GivenName,
		&user.FamilyName,
		&user.Picture,
		&user.Locale,
		&user.VerifiedEmail,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// UpdateUser updates an existing user's information
func (r *usersRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET email = ?, name = ?, given_name = ?, family_name = ?, picture = ?, locale = ?, verified_email = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(
		query,
		user.Email,
		user.Name,
		user.GivenName,
		user.FamilyName,
		user.Picture,
		user.Locale,
		user.VerifiedEmail,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// DeleteUser deletes a user by their ID
func (r *usersRepository) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// CreateSession creates a new session for a user
func (r *usersRepository) CreateSession(session *models.Session) error {
	query := `
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		session.UserID,
		session.Token,
		session.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	session.ID = id
	session.CreatedAt = time.Now()

	return nil
}

// GetSessionByToken retrieves a session by its token
func (r *usersRepository) GetSessionByToken(token string) (*models.Session, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM sessions
		WHERE token = ? AND expires_at > CURRENT_TIMESTAMP
	`

	session := &models.Session{}
	err := r.db.QueryRow(query, token).Scan(
		&session.ID,
		&session.UserID,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session by token: %w", err)
	}

	return session, nil
}

// GetUserBySessionToken retrieves a user by their session token
func (r *usersRepository) GetUserBySessionToken(token string) (*models.User, error) {
	query := `
		SELECT u.id, u.google_id, u.email, u.name, u.given_name, u.family_name, u.picture, u.locale, u.verified_email, u.created_at, u.updated_at
		FROM users u
		INNER JOIN sessions s ON u.id = s.user_id
		WHERE s.token = ? AND s.expires_at > CURRENT_TIMESTAMP
	`

	user := &models.User{}
	err := r.db.QueryRow(query, token).Scan(
		&user.ID,
		&user.GoogleID,
		&user.Email,
		&user.Name,
		&user.GivenName,
		&user.FamilyName,
		&user.Picture,
		&user.Locale,
		&user.VerifiedEmail,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by session token: %w", err)
	}

	return user, nil
}

// DeleteSession deletes a session by its token
func (r *usersRepository) DeleteSession(token string) error {
	query := `DELETE FROM sessions WHERE token = ?`

	result, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

// DeleteExpiredSessions removes all expired sessions from the database
func (r *usersRepository) DeleteExpiredSessions() error {
	query := `DELETE FROM sessions WHERE expires_at <= CURRENT_TIMESTAMP`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}
