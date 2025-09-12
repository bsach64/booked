package httpdelivery

import (
	"net/http"

	httphandler "github.com/bsach64/booked/delivery/http/handler"
	httpmiddleware "github.com/bsach64/booked/delivery/http/middleware"
)

func (s *Server) addRoutes() {
	coreHandler := httphandler.New(s.usecases, s.repositiories)
	middlewares := httpmiddleware.New(s.config, s.usecases, s.repositiories)

	s.serverMux.HandleFunc("/health/", coreHandler.HealthHandler)

	// user
	s.serverMux.HandleFunc("POST /user/register/", coreHandler.RegisterUser)
	s.serverMux.HandleFunc("POST /user/login/", coreHandler.LoginUser)

	// event
	s.serverMux.Handle("POST /event/", middlewares.JWTAuth(middlewares.SetUserInCtx(middlewares.Admin(http.HandlerFunc(coreHandler.CreateEventHandler)))))
	s.serverMux.HandleFunc("GET /event/", coreHandler.GetPaginatedEvents)
	s.serverMux.Handle("DELETE /event/", middlewares.JWTAuth(middlewares.SetUserInCtx(middlewares.Admin(http.HandlerFunc(coreHandler.DeleteEvent)))))
}
