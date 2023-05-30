package controllers

import (
	"github.com/kryptomind/bidboxapi/AccountsService/api/middleware"
)

func (r *Server) initializeRoutes() {
	s := r.Router.PathPrefix("/accounts").Subrouter()
	s.HandleFunc("/", middleware.MiddlewareJSON(r.Home)).Methods("GET")

	//accounts routes
	s.HandleFunc("/user", middleware.ValidateEmail(r.GetUserBalanceByExchange)).Methods("GET")
	s.HandleFunc("/connected", middleware.ValidateEmail(r.GetUserConnectedAccounts)).Methods("GET")
	s.HandleFunc("/exchange-keys", middleware.ValidateEmail(r.GetUserExchangeKeys)).Methods("GET")
	s.HandleFunc("/all-keys", middleware.ValidateEmail(r.GetUserAllKeys)).Methods("GET")
	s.HandleFunc("/disconnect", middleware.ValidateEmail(r.DisconnectKey)).Methods("DELETE")
	s.HandleFunc("/key/settings", middleware.ValidateEmail(r.UpdateKeySettings)).Methods("PUT")
	s.HandleFunc("/key/settings", middleware.ValidateEmail(r.GetSettings)).Methods("GET")

	//trades
	s.HandleFunc("/active_trades", middleware.ValidateEmail(r.GetOpenTrades)).Methods("GET")
	s.HandleFunc("/bitget_history", middleware.ValidateEmail(r.GetClosedTrades)).Methods("GET")
	// s.HandleFunc("/order", middleware.ValidateEmail(r.PlaceOrder)).Methods("POST")
	// s.HandleFunc("/cancel_order", middleware.ValidateEmail(r.DeleteOrder)).Methods("POST")

	s.HandleFunc("/history", middleware.ValidateEmail(r.GetHistory)).Methods("GET")

	//coins
	s.HandleFunc("/admin/coins", middleware.MiddlewareJSON(r.GetCoins)).Methods("GET")
	s.HandleFunc("/admin/coins", middleware.MiddlewareJSON(r.CreateCoin)).Methods("POST")
	s.HandleFunc("/admin/coins", middleware.MiddlewareJSON(r.UpdateCoin)).Methods("PUT")
	s.HandleFunc("/admin/coins", middleware.MiddlewareJSON(r.DeleteCoin)).Methods("DELETE")
	s.HandleFunc("/admin/users", middleware.MiddlewareJSON(r.GetNoOfUsers)).Methods("GET")

	//settings
	s.HandleFunc("/admin/conditions", middleware.MiddlewareJSON(r.GetConditions)).Methods("GET")
	s.HandleFunc("/admin/conditions", middleware.MiddlewareJSON(r.UpdateConditions)).Methods("PUT")

	s.HandleFunc("/register", middleware.MiddlewareJSON(r.SignUpUser)).Methods("POST")
	s.HandleFunc("/login", middleware.MiddlewareJSON(r.LoginUser)).Methods("POST")
	s.HandleFunc("/otp/generate", middleware.MiddlewareJSON(r.GenerateOTP)).Methods("POST")
	s.HandleFunc("/otp/verify", middleware.MiddlewareJSON(r.VerifyOTP)).Methods("POST")
	s.HandleFunc("/otp/validate", middleware.MiddlewareJSON(r.ValidateOTP)).Methods("POST")

}
