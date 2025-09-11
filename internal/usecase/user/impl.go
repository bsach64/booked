package useruc

import (
	"context"

	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/bsach64/booked/utils"
)

type impl struct {
	config   *utils.Config
	userRepo userdom.Repository
}

func (i *impl) CreateUser(ctx context.Context, user userdom.User) error {
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
