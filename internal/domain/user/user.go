package userdom

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type UserRole string

const USER UserRole = "user"
const ADMIN UserRole = "admin"

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email            string `json:"email"`
	UnhashedPassword string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

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
	LoginUser(ctx context.Context, loginRequest *LoginRequest) (*LoginResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
