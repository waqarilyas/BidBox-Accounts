package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

type UserSettingsResp struct {
	Strategy     string `json:"strategy"`
	Mode         string `json:"mode"`
	AutoCompound bool   `json:"auto_compound"`
}

// Get User Settings godoc
// @Summary      Get User Settings
// @Description  Get User Settings
// @Tags         keys
// @Accept       json
// @Produce      json
// @Param        email  query  string true  "email"
// @Success      200  {object}  UserSettingsResp
// @Failure      500  {string}  coin required
// @Router       /keys/settings [get]
func (server *Server) GetSettings(w http.ResponseWriter, r *http.Request, email string) {
	key := models.Key{}
	keys, err := key.GetSettings(server.DB, email)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res := make(map[string]interface{})
	for _, v := range *keys {
		res["strategy"] = v.Strategy
		res["mode"] = v.Mode
		res["auto_compound"] = v.Compound
	}

	response.JSON(w, http.StatusOK, res)
}

type NoUsersResp struct {
	Users        int `json:"users"`
	ActiveTrades int `json:"active_trades"`
	Successful   int `json:"successful"`
	Transactions int `json:"transactions"`
}

// Get Number Of Users godoc
// @Summary      Get Number Of Users
// @Description  Get Number Of Users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  NoUsersResp
// @Failure      500  {string}  coin required
// @Router       /hedge-positions [get]
func (server *Server) GetNoOfUsers(w http.ResponseWriter, r *http.Request) {

	key := models.Key{}
	c, err := key.FindNoOfUsers(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	positions := models.Position{}
	c1, err := positions.GetActivePositions(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	hist := models.Position{}
	c2, err := hist.GetSuccessfulPositions(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res := make(map[string]interface{})
	res["users"] = c
	res["active_trades"] = c1
	res["successful"] = c2
	res["transactions"] = c1 + c2
	response.JSON(w, http.StatusOK, res)
}

// Get Statements Of User godoc
// @Summary      Get Statements Of User
// @Description  Get Statements Of User
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        email  query  string true  "email"
// @Param        service  query  string true  "service"
// @Success      200  {object}  []models.Statements
// @Failure      400  {string}  coin required
// @Failure      500  {string}  coin required
// @Router       /hedge-positions [get]
func (server *Server) GetStatements(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	st := models.Statements{}
	sts, err := st.FindStatements(server.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(*sts) == 0 {
		res := make(map[string]string)
		res["msg"] = "record not found"
		response.JSON(w, http.StatusOK, res)
		return
	}
	response.JSON(w, http.StatusOK, sts)
}

// Get User Trade History godoc
// @Summary      Get User Trade History
// @Description  Get User Trade History
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        email  query  string true  "email"
// @Param        service  query  string true  "service"
// @Success      200  {object}  []models.History
// @Failure      400  {string}  coin required
// @Failure      500  {string}  coin required
// @Router       /hedge-positions [get]
func (server *Server) GetHistory(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	hist := models.History{}
	history, err := hist.GetHistory(server.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(*history) == 0 {
		res := make(map[string]string)
		res["msg"] = "record not found"
		response.JSON(w, http.StatusOK, res)
		return
	}
	response.JSON(w, http.StatusOK, history)
}

// Update Key Settings godoc
// @Summary      Update Key Settings
// @Description  Update Key Settings
// @Tags         keys
// @Accept       json
// @Produce      json
// @Param        email  query  string true  "email"
// @Param        service  query  string true  "service"
// @Param        key  body  models.Key true  "key"
// @Success      200  {object}  []models.Key
// @Failure      400  {string}  coin required
// @Failure      422  {string}  coin required
// @Failure      500  {string}  coin required
// @Router       /keys/settings [put]
func (server *Server) UpdateKeySettings(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	key := models.Key{}
	err = json.Unmarshal(body, &key)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	prev, err := key.FindKeysByEmail(server.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if err := key.Validate(prev); err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	updatedkey, err := key.UpdateKeySettings(server.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, updatedkey)
}

// Disconnect Key godoc
// @Summary      Disconnect Key
// @Description  Disconnect Key
// @Tags         keys
// @Accept       json
// @Produce      json
// @Param        email  query  string true  "email"
// @Param        service  query  string true  "service"
// @Success      200  {string}  res disconnect key
// @Failure      400  {string}  coin required
// @Failure      500  {string}  coin required
// @Router       /disconnect [delete]
func (s *Server) DisconnectKey(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	key := models.Key{}
	_, err := key.DeleteKey(s.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, "disconnected key")
}

// Sync Client Data With Backend godoc
// @Summary      Sync Client Data With Backend
// @Description  Sync Client Data With Backend
// @Tags         keys
// @Accept       json
// @Produce      json
// @Param        user  body  models.User true  "users"
// @Success      200  {object}  models.User
// @Failure      400  {string}  coin required
// @Failure      500  {string}  coin required
// @Router       /sync-account-data [post]
func (s *Server) SyncDataWithClientBackend(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	var user models.User

	err = json.Unmarshal(body, &user)
	if err != nil {
		response.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if user.Email == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("email is required"))
		return
	}

	existingUser, err := user.FindUserByEmail(s.DB, user.Email)

	if err != nil {
		newUser := &models.User{
			Name:     user.Name,
			Country:  user.Country,
			UserName: user.UserName,
			Phone:    user.Phone,
			TimeZone: user.TimeZone,
			Email:    user.Email,
		}
		err := s.DB.Debug().Create(newUser).Error
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.JSON(w, http.StatusOK, newUser)

	} else {
		changes := make(map[string]interface{})

		if user.Name != "" {
			changes["name"] = user.Name
		}
		if user.Country != "" {
			changes["country"] = user.Country
		}
		if user.UserName != "" {
			changes["user_name"] = user.UserName
		}
		if user.Phone != "" {
			changes["phone"] = user.Phone
		}
		if user.TimeZone != "" {
			changes["time_zone"] = user.TimeZone
		}

		if len(changes) > 0 {
			err := s.DB.Debug().Model(&models.User{}).Where("email = ?", user.Email).Updates(changes).Error
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		response.JSON(w, http.StatusOK, existingUser)
		return
	}
}
