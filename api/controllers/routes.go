package controllers

import "github.com/kryptomind/bidboxapi/AccountsService/api/middleware"

func (r *Server) initializeRoutes() {
	s := r.Router.PathPrefix("/accounts").Subrouter()
	s.HandleFunc("/", middleware.MiddlewareJSON(r.Home)).Methods("GET")
	//accounts routes
	s.HandleFunc("/connected", middleware.MiddlewareJSON(r.GetUserConnectedAccounts)).Methods("GET")

}
