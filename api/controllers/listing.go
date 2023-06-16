package controllers

import (
	"errors"
	"net/http"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

func (server *Server) ListToday(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	st := models.Statements{}
	sts, err := st.GetStatementsToday(server.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(*sts) == 0 {
		res := make(map[string]string)
		res["msg"] = "no statements today"
		response.JSON(w, http.StatusOK, res)
		return
	}
	response.JSON(w, http.StatusOK, sts)
}

func (server *Server) ListTotal(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	st := models.Statements{}
	sts, err := st.GetStatementsAllTime(server.DB, email, service)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	if len(*sts) == 0 {
		res := make(map[string]string)
		res["msg"] = "no statements for this user"
		response.JSON(w, http.StatusOK, res)
		return
	}
	response.JSON(w, http.StatusOK, sts)
}

func (server *Server) ListCoinwise(w http.ResponseWriter, r *http.Request, email string) {
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
