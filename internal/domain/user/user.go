package userdom

import (
	"context"
)

type UserRole string

const USER UserRole = "user"
const ADMIN UserRole = "admin"

type NewUser struct {
	Name             string `json:"name"`
	UnhashedPassword string `json:"password"`
	Email            string `json:"email"`
}

type User struct {
	Name           string
	HashedPassword string
	Email          string
	Role           UserRole
}

type Usecase interface {
	CreateUser(ctx context.Context, user NewUser) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
