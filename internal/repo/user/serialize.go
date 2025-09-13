package userrepo

import (
	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/bsach64/booked/internal/repo/sql/db"
)

func ToUserDomain(user db.User) *userdom.User {
	return &userdom.User{
		ID:             user.ID.Bytes,
		Name:           user.Name,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Role:           userdom.UserRole(user.Role),
	}
}
