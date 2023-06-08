package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

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

func (server *Server) GetNoOfUsers(w http.ResponseWriter, r *http.Request) {

	key := models.Key{}
	c, err := key.FindNoOfUsers(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	order := models.Order{}
	c1, err := order.GetActiveTrades(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	hist := models.History{}
	c2, err := hist.GetSuccessfulTrades(server.DB)
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

	fmt.Println("---previous---", prev)

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
