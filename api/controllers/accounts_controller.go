package controllers

import (
	"errors"
	"net/http"

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
type ExchangeResponse struct {
	Name      string `json:"name"`
	Short     string `json:"short"`
	ImageSrc  string `json:"image_src"`
	Id        int    `json:"id"`
	Connected bool   `json:"connected"`
	IsActive  bool   `json:"is_active"`
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

	//get all exchanges
	exchange := models.Exchanges{}
	exchanges, err := exchange.FindAllExchanges(server.DB)
	if err != nil {
		// if there are no connected accounts, return a default response
		response.ERROR(w, http.StatusBadRequest, errors.New("no exchanges Found"))
		return
	}

	exchangeResponses := make([]ExchangeResponse, len(*exchanges))

	for i, e := range *exchanges {
		exchangeResponses[i] = ExchangeResponse{
			Name:      e.Name,
			Short:     e.Short,
			ImageSrc:  e.ImageSrc,
			Id:        int(e.Id.ID()),
			Connected: false,
			IsActive:  e.IsActive,
		}
	}

	//get user by email address
	// user := models.User{}
	// userRes, err := user.FindUserByEmail(server.DB, strings.ToLower(email))
	// if err != nil {
	// 	// response.ERROR(w, http.StatusBadRequest, errors.New("no user found with given email address"))
	// 	response.JSON(w, http.StatusOK, exchangeResponses)

	// 	return
	// }

	//get keys by user id
	key := models.Key{}
	keys, err := key.FindKeysByUserEmail(server.DB, email)
	if err != nil {
		response.JSON(w, http.StatusOK, exchangeResponses)
		return
	}

	for _, value := range *keys {
		for j, exchangeResponse := range exchangeResponses {
			if exchangeResponse.Short == value.Service {
				exchangeResponses[j].Connected = true
				break
			}
		}
	}

	// return response as JSON
	response.JSON(w, http.StatusOK, exchangeResponses)
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
	// user := models.User{}
	// userRes, err := user.FindUserByEmail(server.DB, strings.ToLower(email))
	// if err != nil {
	// 	response.ERROR(w, http.StatusBadRequest, errors.New("no user found with given email address"))
	// 	return
	// }

	//get all exchanges
	exchangeMod := models.Exchanges{}
	dbExchange, err := exchangeMod.GetExchangeByShort(server.DB, exchange)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("invalid exchange id"))
		return
	}

	//get keys by user id
	key := models.Key{}
	keys, err := key.FindKeyByUserEmailAndShort(server.DB, email, dbExchange.Short)
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

func (server *Server) UserCheck(w http.ResponseWriter, r *http.Request) {
	exchange := models.Exchanges{}
	exchanges, err := exchange.FindAllExchanges(server.DB)
	if err != nil {
		// if there are no connected accounts, return a default response
		response.ERROR(w, http.StatusBadRequest, errors.New("no exchanges Found"))
		return
	}

	response.JSON(w, http.StatusOK, exchanges)

}
