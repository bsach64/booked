package httpdelivery

import (
	httphandler "github.com/bsach64/booked/delivery/http/handler"
	httpmiddleware "github.com/bsach64/booked/delivery/http/middleware"
)

func (s *Server) addRoutes() {
	coreHandler := httphandler.New(s.usecases, s.repositiories)
	_ = httpmiddleware.New(s.config, s.usecases, s.repositiories)

	s.serverMux.HandleFunc("/health/", coreHandler.HealthHandler)

	// user
	s.serverMux.HandleFunc("POST /user/register/", coreHandler.RegisterUser)
	s.serverMux.HandleFunc("POST /user/login/", coreHandler.LoginUser)
}
