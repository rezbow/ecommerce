package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	IsAdmin      bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserSvc interface {
	RegisterUser(*RegisterUser) (*User, error)
	Authenticate(*Login) (string, error)
	GetUser(uuid.UUID) (*User, error)
}

type UserRepo interface {
	Get(string) (*User, error)
	GetByEmail(string) (*User, error)
	Create(*User) error
}
