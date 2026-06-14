package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"social-network/pkg/apperror"
	"social-network/pkg/models"
	"strings"
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

// func NewUserRepository(db *sql.DB) UserRepository {
// 	return &sqliteUserRepo{db: db}
// }

func (r *sqliteUserRepo) CreateUser(u *models.User) error {
	const query = `
		INSERT INTO users (id, email, password, first_name, last_name, dob, nickname, 
		about_me, avatar_path, is_public, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, 
		u.ID, u.Email, u.Password, u.FirstName, u.LastName,
		u.DOB, u.Nickname, u.AboutMe, u.AvatarPath,
		booleanToInt(u.IsPublic), u.CreatedAt,
	)
	if err != nil {
		if isSQLiteUniqueViolation(err) {
			return apperror.Conflict("email already registered")
		}

		// TODO:
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
	var isPublic int

	err := r.db.QueryRow(query, id).Scan(
		&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.DOB,
		&u.Nickname, &u.AboutMe, &u.AvatarPath, &isPublic, &u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound("user not found")
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	u.IsPublic = isPublic == 1
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

func booleanToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func isSQLiteUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()

	return strings.Contains(msg, "UNIQUE constraint failed") || strings.Contains(msg, "unique constraint failed")
}

// requireOneRow returns an error if the statement affected no rows.
// This is because SQLite treats updating/deleting a row that doesn't exist (0 rows affected) as a success
func requireOneRow(res sql.Result, entity string) error {
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return apperror.NotFound(entity + " not found")
	}
	return nil
}