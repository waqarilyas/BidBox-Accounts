package controllers

import "github.com/kryptomind/bidboxapi/AccountsService/api/middleware"

func (r *Server) initializeRoutes() {
	s := r.Router.PathPrefix("/accounts").Subrouter()
	s.HandleFunc("/", middleware.MiddlewareJSON(r.Home)).Methods("GET")

	//accounts routes
	s.HandleFunc("/user", middleware.MiddlewareJSON(r.GetUserAcounts)).Methods("POST")

	//trades
	s.HandleFunc("/active_trades", middleware.MiddlewareJSON(r.GetOpenTrades)).Methods("GET")
	s.HandleFunc("/history", middleware.MiddlewareJSON(r.GetClosedTrades)).Methods("GET")
	s.HandleFunc("/connected", middleware.MiddlewareJSON(r.GetUserConnectedAccounts)).Methods("GET")

}
