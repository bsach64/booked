package httpmiddleware

import (
	"github.com/bsach64/booked/internal/repo"
	"github.com/bsach64/booked/internal/usecase"
	"github.com/bsach64/booked/utils"
)

type CoreMiddleware struct {
	config       *utils.Config
	usecases     usecase.Usecase
	repositories repo.Repositories
}

func New(config *utils.Config, usecases usecase.Usecase, repositories repo.Repositories) *CoreMiddleware {
	return &CoreMiddleware{
		config:       config,
		usecases:     usecases,
		repositories: repositories,
	}
}
