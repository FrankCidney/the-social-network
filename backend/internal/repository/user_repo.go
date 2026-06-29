package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/internal/apperror"
	"social-network/internal/models"
)

type UserRepository interface {
	CreateUser(u *models.User) error
	GetUserByID(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(u *models.User) error
	SetProfileVisibility(userID string, isPublic bool) error
	UpdateAvatarPath(userID, path string) error
}

type sqliteUserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &sqliteUserRepo{db: db}
}

func (r *sqliteUserRepo) CreateUser(u *models.User) error {
	const query = `
		INSERT INTO users (id, email, password, first_name, last_name, dob, nickname, 
		about_me, avatar_path, is_public, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query,
		u.ID, u.Email, u.Password, u.FirstName, u.LastName,
		u.DOB, u.Nickname, u.AboutMe, u.AvatarPath,
		u.IsPublic, u.CreatedAt,
	)
	if err != nil {
		if isSQLiteUniqueViolation(err) {
			return apperror.Conflict("email already registered")
		}

		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *sqliteUserRepo) GetUserByID(id string) (*models.User, error) {
	const query = `
		SELECT id, email, password, first_name, last_name, dob,
			COALESCE(nickname, ''), COALESCE(about_me, ''), COALESCE(avatar_path, ''),
			is_public, created_at
		FROM users WHERE id = ?`

	u := &models.User{}

	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.DOB,
		&u.Nickname, &u.AboutMe, &u.AvatarPath, &u.IsPublic, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound("user not found")
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return u, nil
}

func (r *sqliteUserRepo) UpdateUser(u *models.User) error {
	const query = `
		UPDATE users
		SET first_name = ?, last_name = ?, dob = ?, nickname = ?, about_me = ?
		WHERE id = ?`

	res, err := r.db.Exec(query,
		u.FirstName, u.LastName, u.DOB, u.Nickname, u.AboutMe, u.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return requireOneRow(res, "user")
}

func (r *sqliteUserRepo) GetUserByEmail(email string) (*models.User, error) {
	const query = `
		SELECT id, email, password, first_name, last_name, dob,
		       COALESCE(nickname, ''), COALESCE(about_me, ''), COALESCE(avatar_path, ''),
		       is_public, created_at
		FROM users WHERE email = ?`

	u := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.DOB,
		&u.Nickname, &u.AboutMe, &u.AvatarPath, &u.IsPublic, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.Unauthorized("invalid credentials")
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return u, nil
}

func (r *sqliteUserRepo) SetProfileVisibility(userID string, isPublic bool) error {
	const query = `UPDATE users SET is_public = ? WHERE id = ?`
	res, err := r.db.Exec(query, isPublic, userID)
	if err != nil {
		return fmt.Errorf("set visibility: %w", err)
	}
	return requireOneRow(res, "user")
}

func (r *sqliteUserRepo) UpdateAvatarPath(userID, path string) error {
	const query = `UPDATE users SET avatar_path = ? WHERE id = ?`
	res, err := r.db.Exec(query, path, userID)
	if err != nil {
		return fmt.Errorf("update avatar: %w", err)
	}
	return requireOneRow(res, "user")
}
