package controllers

import (
	"net/http"

	swag "github.com/go-openapi/runtime/middleware"
	"github.com/kryptomind/bidboxapi/AccountsService/api/middleware"
)

func (r *Server) initializeRoutes() {

	r.Router.Handle("/", middleware.MiddlewareJSON(r.Home)).Methods("GET")

	r.Router.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// documentation for developers
	opts := swag.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	sh := swag.SwaggerUI(opts, nil)
	r.Router.Handle("/docs", sh)

	// documentation for share
	// opts1 := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	// sh1 := middleware.Redoc(opts1, nil)
	// r.Handle("/docs", sh1)

	s := r.Router.PathPrefix("/accounts").Subrouter()

	//accounts routes
	s.HandleFunc("/user", middleware.ValidateEmail(r.GetUserBalanceByExchange)).Methods("GET")
	s.HandleFunc("/connected", middleware.ValidateEmail(r.GetUserConnectedAccounts)).Methods("GET")
	s.HandleFunc("/exchange-keys", middleware.ValidateEmail(r.GetUserExchangeKeys)).Methods("GET")
	s.HandleFunc("/all-keys", middleware.ValidateEmail(r.GetUserAllKeys)).Methods("GET")
	s.HandleFunc("/disconnect", middleware.ValidateEmail(r.DisconnectKey)).Methods("DELETE")
	s.HandleFunc("/key/settings", middleware.ValidateEmail(r.UpdateKeySettings)).Methods("PUT")
	s.HandleFunc("/key/settings", middleware.ValidateEmail(r.GetSettings)).Methods("GET")
	s.HandleFunc("/sync-account-data", r.SyncDataWithClientBackend).Methods("POST")

	//trades
	s.HandleFunc("/active_trades", middleware.ValidateEmail(r.GetOpenTrades)).Methods("GET")
	s.HandleFunc("/bitget_history", middleware.ValidateEmail(r.GetClosedTrades)).Methods("GET")
	// s.HandleFunc("/order", middleware.ValidateEmail(r.PlaceOrder)).Methods("POST")
	// s.HandleFunc("/cancel_order", middleware.ValidateEmail(r.DeleteOrder)).Methods("POST")
	s.HandleFunc("/statements", middleware.ValidateEmail(r.GetStatements)).Methods("GET")
	s.HandleFunc("/leaderboard", middleware.MiddlewareJSON(r.GetLeaderBoardv2)).Methods("GET")

	//coins
	s.HandleFunc("/admin/coins", middleware.MiddlewareJSON(r.GetCoins)).Methods("GET")
	s.HandleFunc("/admin/coins", middleware.MiddlewareJSON(r.CreateCoin)).Methods("POST")
	s.HandleFunc("/admin/coins", middleware.MiddlewareJSON(r.UpdateCoin)).Methods("PUT")
	s.HandleFunc("/admin/coins", middleware.MiddlewareJSON(r.DeleteCoin)).Methods("DELETE")
	s.HandleFunc("/admin/stats", middleware.MiddlewareJSON(r.GetNoOfUsers)).Methods("GET")

	//settings
	s.HandleFunc("/admin/conditions", middleware.MiddlewareJSON(r.GetConditions)).Methods("GET")
	s.HandleFunc("/admin/conditions", middleware.MiddlewareJSON(r.UpdateConditions)).Methods("PUT")

	// 2FA
	s.HandleFunc("/register", middleware.MiddlewareJSON(r.SignUpUser)).Methods("POST")
	s.HandleFunc("/login", middleware.MiddlewareJSON(r.LoginUser)).Methods("POST")
	s.HandleFunc("/otp/enable", middleware.MiddlewareJSON(r.EnableOTP)).Methods("POST")
	s.HandleFunc("/otp/generate", middleware.MiddlewareJSON(r.GenerateOTP)).Methods("POST")
	s.HandleFunc("/otp/validate", middleware.MiddlewareJSON(r.ValidateOTP)).Methods("POST")
	s.HandleFunc("/otp/verify", middleware.MiddlewareJSON(r.VerifyOTP)).Methods("POST")
	s.HandleFunc("/changePassword", middleware.MiddlewareJSON(r.ChangePassword)).Methods("POST")

	// admin settings
	s.HandleFunc("/changeTimeframe", middleware.MiddlewareJSON(r.changeTimeframe)).Methods("POST")
	s.HandleFunc("/changeSettings", middleware.MiddlewareJSON(r.changeSettings)).Methods("POST")
	s.HandleFunc("/getSettings", middleware.MiddlewareJSON(r.getSettings)).Methods("GET")

	s.HandleFunc("/getTimeframe", middleware.MiddlewareJSON(r.GetTimeframe)).Methods("GET")
	s.HandleFunc("/history", middleware.ValidateEmail(r.GetPositionHistory)).Methods("GET")
	// s.HandleFunc("/closePositions", middleware.MiddlewareJSON(r.GetClosedPositionsByEmail)).Methods("GET")

	// listing
	s.HandleFunc("/listing/today", middleware.ValidateEmail(r.ListToday)).Methods("GET")
	s.HandleFunc("/listing/total", middleware.ValidateEmail(r.ListTotal)).Methods("GET")
	s.HandleFunc("/listing/coinwise", middleware.ValidateEmail(r.ListCoinwise)).Methods("GET")

	// updated code afetr hedging logic

	s.HandleFunc("/orders", middleware.ValidateEmail(r.GetUserOrders)).Methods("GET")
	s.HandleFunc("/positions", middleware.ValidateEmail(r.GetUserOpenPositions)).Methods("GET")
	s.HandleFunc("/hedge-positions", middleware.ValidateEmail(r.UserClearedHedgePositions)).Methods("GET")

}
