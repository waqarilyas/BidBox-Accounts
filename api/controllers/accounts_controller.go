package controllers

import (
	"errors"
	"net/http"

	"github.com/kryptomind/bidboxapi/AccountsService/api/helpers"
	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "Accounts Service")
}

type UserConnectedAccountsRequest struct {
	Email string `json:"email"`
}

// An ExchangeModel is gives information about the connected exchanges
// swagger:model ExchangeModel
type ExchangeModel struct {
	// Example: bybit
	Name string `json:"name"`
	// Example: bybit
	Short string `json:"short"`
	// Example: bybit
	ImageSrc string `json:"image_src"`
	// Example: 23
	Id int `json:"id"`
	// Example: true
	Connected bool `json:"connected"`
	// Example: true
	IsActive bool `json:"is_active"`
}

// swagger:model ExchangeRes
type ExchangeRes struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/ExchangeModel"
	//  required: true
	Body ExchangeModel `json:"body"`
}

// swagger:route GET /connected pets users listPets
//
// Get Connected exchanges for a user.
//
// This will show all available pets by default.
// You can get the pets that are out of stock
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Parameters:
//       + name: email
//         in: query
//         description: user email
//         required: true
//         type: string
//
//     Responses:
//       default: genericError
//       200: ExchangeRes
//       404: connected account not found
func (server *Server) GetUserConnectedAccounts(w http.ResponseWriter, r *http.Request, email string) {

	//get all exchanges
	exchange := models.Exchanges{}
	exchanges, err := exchange.FindAllExchanges(server.DB)
	if err != nil {
		// if there are no connected accounts, return a default response
		response.ERROR(w, http.StatusNotFound, errors.New("no exchanges Found"))
		return
	}

	exchangeResponses := make([]ExchangeModel, len(*exchanges))

	for i, e := range *exchanges {
		exchangeResponses[i] = ExchangeModel{
			Name:      e.Name,
			Short:     e.Short,
			ImageSrc:  e.ImageSrc,
			Id:        int(e.Id.ID()),
			Connected: false,
			IsActive:  e.IsActive,
		}
	}

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

func (server *Server) GetUserBalanceByExchange(w http.ResponseWriter, r *http.Request, email string) {
	exchange := r.URL.Query().Get("exchange")

	if exchange == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("exchange is required"))
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

func (server *Server) GetUserAllKeys(w http.ResponseWriter, r *http.Request, email string) {
	userEmail := r.URL.Query().Get("email")
	if userEmail == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("user email is required"))
		return
	}

	key := models.Key{}
	apiKeys, err := key.FindKeysByUserEmail(server.DB, userEmail)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	keys := *apiKeys

	for i, v := range *apiKeys {

		keys[i].ApiKey, err = helpers.DecryptStrings(v.ApiKey)
		if err != nil {
			response.ERROR(w, http.StatusBadRequest, errors.New("something went wrong while decrypting api key"))
			return
		}

		if v.Passphrase != "" {
			keys[i].Passphrase, err = helpers.DecryptStrings(v.Passphrase)
			if err != nil {
				response.ERROR(w, http.StatusBadRequest, errors.New("something went wrong while decrypting passphrase"))
				return
			}
		}
		keys[i].SecretKey, err = helpers.DecryptStrings(v.SecretKey)
		if err != nil {
			response.ERROR(w, http.StatusBadRequest, errors.New("something went wrong while decrypting secret key"))
			return
		}

	}

	response.JSON(w, http.StatusOK, keys)
}

func (server *Server) GetUserExchangeKeys(w http.ResponseWriter, r *http.Request, email string) {
	exchange := r.URL.Query().Get("exchange")
	if exchange == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("exchange is required"))
		return
	}

	userEmail := r.URL.Query().Get("email")
	if exchange == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("user email is required"))
		return
	}

	key := models.Key{}
	apiKeys, err := key.FindKeyByUserEmailAndShort(server.DB, userEmail, exchange)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("provided exchange not connected"))
		return
	}

	decrypted_api_key, err := helpers.DecryptStrings(apiKeys.ApiKey)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("something went wrong while getting api key"))
		return

	}
	decrypted_passphrase, err := helpers.DecryptStrings(apiKeys.Passphrase)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("something went wrong while getting api key"))
		return

	}
	decrypted_secret, err := helpers.DecryptStrings(apiKeys.SecretKey)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, errors.New("something went wrong while getting api key"))
		return
	}

	exchangeResponse := map[string]string{
		"api_key":    decrypted_api_key,
		"passphrase": decrypted_passphrase,
		"secret":     decrypted_secret,
	}

	response.JSON(w, http.StatusOK, exchangeResponse)
}
