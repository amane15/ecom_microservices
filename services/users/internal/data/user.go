package data

import (
	"database/sql"
	"time"
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
	return nil, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	return nil, nil
}

func (m UserModel) Create(user *User) error {
	return nil
}
