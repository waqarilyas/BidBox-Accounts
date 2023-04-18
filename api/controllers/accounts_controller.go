package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "Accounts Service")
}

type UserConnectedAccountsRequest struct {
	Email string `json:"email"`
}

func (server *Server) GetUserConnectedAccounts(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("email is required"))
		return
	}

	if !govalidator.IsEmail(email) {
		response.ERROR(w, http.StatusBadRequest, errors.New("invalid email address"))
		return
	}

	//get user by email address
	user := models.User{}
	userRes, err := user.FindUserByEmail(server.DB, strings.ToLower(email))
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("no user found with given email address"))
		return
	}

	//get all exchanges
	exchange := models.Exchanges{}
	exchanges, err := exchange.FindAllExchanges(server.DB)
	if err != nil {
		// if there are no connected accounts, return a default response
		response.ERROR(w, http.StatusBadRequest, errors.New("no exchanges Found"))
		return
	}

	responseMap := make(map[string]string)
	for _, e := range *exchanges {
		responseMap[e.Short] = "Not Connected"
	}

	//get keys by user id
	key := models.Key{}
	keys, err := key.FindKeysByUserId(server.DB, userRes.Id)
	if err != nil {
		response.JSON(w, http.StatusOK, responseMap)
		return
	}

	for _, value := range *keys {
		responseMap[value.Service] = "Connected"
	}

	// return response as JSON
	response.JSON(w, http.StatusOK, responseMap)
}

func (server *Server) GetUserBalanceByExchange(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	exchange := r.URL.Query().Get("exchange")

	if email == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("email is required"))
		return
	}
	if exchange == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("exchange is required"))
		return
	}

	if !govalidator.IsEmail(email) {
		response.ERROR(w, http.StatusBadRequest, errors.New("invalid email address"))
		return
	}

	//get user by email address
	user := models.User{}
	userRes, err := user.FindUserByEmail(server.DB, strings.ToLower(email))
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("no user found with given email address"))
		return
	}

	//get all exchanges
	exchangeMod := models.Exchanges{}
	dbExchange, err := exchangeMod.GetExchangeByShort(server.DB, exchange)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("invalid exchange id"))
		return
	}

	//get keys by user id
	key := models.Key{}
	keys, err := key.FindKeyByUserIdAndShort(server.DB, userRes.Id, dbExchange.Short)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("provided exchange not connected"))
		return
	}

	//get account based on key

	account := models.Accounts{}

	dbAccount, err := account.GetAccountByApiKeyId(server.DB, keys.Keyid)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("no account found"))
		return
	}

	response.JSON(w, http.StatusOK, dbAccount)

}

func (server *Server) GetSupportedExchange(w http.ResponseWriter, r *http.Request) {
	exchange := models.Exchanges{}
	exchanges, err := exchange.FindAllExchanges(server.DB)
	if err != nil {
		// if there are no connected accounts, return a default response
		response.ERROR(w, http.StatusBadRequest, errors.New("no exchanges Found"))
		return
	}

	response.JSON(w, http.StatusOK, exchanges)

}
