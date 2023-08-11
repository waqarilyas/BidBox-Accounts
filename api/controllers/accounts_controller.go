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

type ExchangeModel struct {
	Name      string `json:"name"`
	Short     string `json:"short"`
	ImageSrc  string `json:"image_src"`
	Id        int    `json:"id"`
	Connected bool   `json:"connected"`
	IsActive  bool   `json:"is_active"`
}

// Get User Connected Accounts godoc
// @Summary      Get User Connected Accounts
// @Description  Get User Connected Accounts
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        email  query  string true  "email"
// @Success      200  {object}  []ExchangeModel
// @Failure      404  {string}  coin required
// @Router       /connected [get]
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

// Get All User Keys godoc
// @Summary      Get All User Keys
// @Description  Get All User Keys
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        email  query  string true  "email"
// @Success      200  {object}  []models.Key
// @Failure      400  {string}  coin required
// @Failure      500  {string}  coin required
// @Router       /all-keys [get]
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

type ExchangeResp struct {
	ApiKey     string `json:"api_key"`
	Passphrase string `json:"passphrase"`
	Secret     string `json:"secret"`
}

// Get User Exchnage Keys godoc
// @Summary      Get User Exchnage Keys
// @Description  Get User Exchnage Keys
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        email  query  string true  "email"
// @Param        exchange  query  string true  "exchange"
// @Success      200  {object}  ExchangeResp
// @Failure      400  {string}  coin required
// @Router       /exchange-keys [get]
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
