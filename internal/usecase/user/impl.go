package useruc

import (
	"context"

	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/bsach64/booked/utils"
	"golang.org/x/crypto/bcrypt"
)

type impl struct {
	config   *utils.Config
	userRepo userdom.Repository
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func (i *impl) CreateUser(ctx context.Context, NewUser userdom.NewUser) error {
	hashedPassword, err := HashPassword(NewUser.UnhashedPassword)
	if err != nil {
		return err
	}
	user := userdom.User{
		Name:           NewUser.Name,
		Email:          NewUser.Email,
		HashedPassword: hashedPassword,
		Role:           userdom.USER,
	}
	return i.userRepo.CreateUser(ctx, user)
}

func (i *impl) GetUserByEmail(ctx context.Context, email string) (*userdom.User, error) {
	return i.userRepo.GetUserByEmail(ctx, email)
}

func New(config *utils.Config, userRepo userdom.Repository) userdom.Usecase {
	return &impl{
		config:   config,
		userRepo: userRepo,
	}
}
