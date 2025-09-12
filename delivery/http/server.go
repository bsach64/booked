package httpdelivery

import (
	"log"
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

func New(config *utils.Config, usecases usecase.Usecase, repositories repo.Repositories) *Server {
	server := &Server{
		serverMux:     http.NewServeMux(),
		config:        config,
		usecases:      usecases,
		repositiories: repositories,
	}
	server.addRoutes()
	return server
}

func (s *Server) StartServer() {
	err := http.ListenAndServe(s.config.ServerURL, s.serverMux)
	if err != nil {
		log.Fatalf("got err=%v", err)
	}
}
