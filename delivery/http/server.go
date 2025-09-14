package httpdelivery

import (
	"net/http"

	"github.com/bsach64/booked/internal/repo"
	"github.com/bsach64/booked/internal/usecase"
	"github.com/bsach64/booked/utils"
)

type Server struct {
	serverMux     *http.ServeMux
	config        *utils.Config
	usecases      usecase.Usecase
	repositiories repo.Repositories
}

func New(config *utils.Config, usecases usecase.Usecase, repositories repo.Repositories) *http.Server {
	server := &Server{
		serverMux:     http.NewServeMux(),
		config:        config,
		usecases:      usecases,
		repositiories: repositories,
	}
	server.addRoutes()
	return &http.Server{Addr: ":"+config.ServerURL}
}
