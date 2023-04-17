package controllers

import "github.com/kryptomind/bidboxapi/AccountsService/api/middleware"

func (r *Server) initializeRoutes() {
	s := r.Router.PathPrefix("/accounts").Subrouter()
	s.HandleFunc("/", middleware.MiddlewareJSON(r.Home)).Methods("GET")

	//accounts routes
	s.HandleFunc("/user", middleware.MiddlewareJSON(r.GetUserAcounts)).Methods("POST")



}
