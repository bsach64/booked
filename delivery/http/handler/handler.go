package httphandler

import (
	"github.com/bsach64/booked/internal/repo"
	"github.com/bsach64/booked/internal/usecase"
)

type CoreHandler struct {
	usecases     usecase.Usecase
	repositories repo.Repositories
}

func New(usecase usecase.Usecase, repositories repo.Repositories) *CoreHandler {
	return &CoreHandler{
		usecases:     usecase,
		repositories: repositories,
	}
}
