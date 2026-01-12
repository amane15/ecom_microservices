package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/amane15/ecom_microservice/services/user/internal/validator"
	"github.com/lib/pq"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) GetByID(id int) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, email, first_name, last_name, role, created_at, updated_at
	FROM users 
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := &User{}
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	return nil, nil
}

func (m UserModel) Insert(user *User) error {
	query := `INSERT INTO users(email, first_name, last_name)
	VALUES ($1, $2, $3)
	RETURNING role, created_at, updated_at;`

	args := []any{user.Email, user.FirstName, user.LastName}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return ErrDuplicateEmail
			}
		} else {
			return err
		}
	}
	return nil
}

func (m UserModel) Update(user *User) error {
	query := `UPDATE users
	SET first_name = $1, last_name = $2 
	WHERE id = $3
	RETURNING first_name, last_name`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.ID).Scan(
		&user.FirstName,
		&user.LastName,
	)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FirstName != "", "first_name", "must be provided")
	v.Check(len(user.FirstName) <= 255, "first_name", "must not be more thann 255 bytes long")
	v.Check(user.LastName != "", "last_name", "must be provided")
	v.Check(len(user.LastName) <= 255, "last_name", "must not be more thann 255 bytes long")

	v.Check(user.Email != "", "email", "must be provided")
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be a valid email address")

	// v.Check(user.Role != "user" && user.Role != "admin", "role", "must be a valid role")
}
