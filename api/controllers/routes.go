package controllers

import "github.com/kryptomind/bidboxapi/AccountsService/api/middleware"

func (r *Server) initializeRoutes() {

	s := r.Router.PathPrefix("/accounts").Subrouter()
	s.HandleFunc("/", middleware.MiddlewareJSON(r.Home)).Methods("GET")

	//accounts routes
	s.HandleFunc("/user", middleware.ValidateEmail(r.GetUserBalanceByExchange)).Methods("GET")
	s.HandleFunc("/connected", middleware.ValidateEmail(r.GetUserConnectedAccounts)).Methods("GET")

	//trades
	s.HandleFunc("/active_trades", middleware.ValidateEmail(r.GetOpenTrades)).Methods("GET")
	s.HandleFunc("/history", middleware.ValidateEmail(r.GetClosedTrades)).Methods("GET")
	s.HandleFunc("/order", middleware.ValidateEmail(r.PlaceOrder)).Methods("POST")

}
