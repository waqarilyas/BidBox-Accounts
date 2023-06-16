package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

func (server *Server) ListToday(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	// Parse query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	// Calculate offset based on page and limit
	offset := (page - 1) * limit

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	st := models.Statements{}
	sts, err := st.GetStatementsToday(server.DB, email, service, limit, offset)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res := make(map[string]interface{})
	if len(*sts) == 0 {
		res["msg"] = "no statements for this user"
		res["data"] = []string{}
		response.JSON(w, http.StatusOK, res)
		return
	}
	res["msg"] = "success"
	res["data"] = sts
	response.JSON(w, http.StatusOK, res)
}

func (server *Server) ListTotal(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}
	// Parse query parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	// Calculate offset based on page and limit
	offset := (page - 1) * limit

	st := models.Statements{}
	sts, err := st.GetStatementsAllTime(server.DB, email, service, limit, offset)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res := make(map[string]interface{})
	if len(*sts) == 0 {
		res["msg"] = "no statements for this user"
		res["data"] = []string{}
		response.JSON(w, http.StatusOK, res)
		return
	}
	res["msg"] = "success"
	res["data"] = sts
	response.JSON(w, http.StatusOK, res)
}

func (server *Server) ListCoinwise(w http.ResponseWriter, r *http.Request, email string) {
	service := r.URL.Query().Get("service")

	if service == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("service is required"))
		return
	}

	time := r.URL.Query().Get("time")

	st := models.Statements{}
	var sts *[]models.Result
	var err error

	if time == "day" {
		sts, err = st.GetCoinwiseToday(server.DB, email, service)
		if err != nil {
			response.ERROR(w, http.StatusInternalServerError, err)
			return
		}

	} else if time == "month" {
		sts, err = st.GetCoinwiseMonth(server.DB, email, service)
		if err != nil {
			response.ERROR(w, http.StatusInternalServerError, err)
			return
		}

	} else {
		sts, err = st.GetCoinwiseAllTime(server.DB, email, service)
		if err != nil {
			response.ERROR(w, http.StatusInternalServerError, err)
			return
		}

	}

	res := make(map[string]interface{})
	if len(*sts) == 0 {
		res["msg"] = "no statements for this user"
		res["data"] = []string{}
		response.JSON(w, http.StatusOK, res)
		return
	}
	res["msg"] = "success"
	res["data"] = sts
	response.JSON(w, http.StatusOK, res)
}
