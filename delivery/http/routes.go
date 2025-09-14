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
	s.serverMux.Handle("GET /user/bookings/", middlewares.JWTAuth(middlewares.SetUserInCtx(http.HandlerFunc(coreHandler.UserBookings))))

	// event
	s.serverMux.Handle("POST /event/", middlewares.JWTAuth(middlewares.SetUserInCtx(middlewares.Admin(http.HandlerFunc(coreHandler.CreateEventHandler)))))
	s.serverMux.HandleFunc("GET /event/", coreHandler.GetPaginatedEvents)
	s.serverMux.Handle("DELETE /event/", middlewares.JWTAuth(middlewares.SetUserInCtx(middlewares.Admin(http.HandlerFunc(coreHandler.DeleteEvent)))))
	s.serverMux.Handle("POST /event/update/", middlewares.JWTAuth(middlewares.SetUserInCtx(middlewares.Admin(http.HandlerFunc(coreHandler.UpdateEvent)))))

	// tickets
	s.serverMux.Handle("POST /ticket/reserve/", middlewares.JWTAuth(middlewares.SetUserInCtx(http.HandlerFunc(coreHandler.ReserveTickets))))
	s.serverMux.Handle("POST /ticket/book/", middlewares.JWTAuth(middlewares.SetUserInCtx(http.HandlerFunc(coreHandler.BookTickets))))
	s.serverMux.Handle("POST /ticket/cancel/", middlewares.JWTAuth(middlewares.SetUserInCtx(http.HandlerFunc(coreHandler.CancelTickets))))

	// analytics
	s.serverMux.Handle("GET /analytics/", middlewares.JWTAuth(middlewares.SetUserInCtx(middlewares.Admin(http.HandlerFunc(coreHandler.TotalBookingsHandler)))))
	s.serverMux.Handle("GET /analytics/cancellation_rates/", middlewares.JWTAuth(middlewares.SetUserInCtx(middlewares.Admin(http.HandlerFunc(coreHandler.CancellationRatesHandler)))))

	// waitlist
	s.serverMux.Handle("POST /waitlist/add/", middlewares.JWTAuth(middlewares.SetUserInCtx(http.HandlerFunc(coreHandler.AddWaitlistHandler))))
}
