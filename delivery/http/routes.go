package httpdelivery

import (
	httphandler "github.com/bsach64/booked/delivery/http/handler"
)

func (s *Server) addRoutes() {
	s.serverMux.HandleFunc("/health", httphandler.HealthHandler)
}
