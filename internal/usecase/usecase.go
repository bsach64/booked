package usecase

import (
	userdom "github.com/bsach64/booked/internal/domain/user"
	"github.com/bsach64/booked/internal/repo"
	useruc "github.com/bsach64/booked/internal/usecase/user"
	"github.com/bsach64/booked/utils"
)

type Usecase struct {
	config *utils.Config
	UserUC userdom.Usecase
}

func New(config *utils.Config, repositories repo.Repositories) Usecase {
	return Usecase{
		config: config,
		UserUC: useruc.New(config, repositories.User),
	}
}
