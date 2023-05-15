package controllers

import (
	"net/http"
	"strconv"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models/admin"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

func (server *Server) GetConditions(w http.ResponseWriter, r *http.Request) {

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

	cond := admin.Conditions{}

	conds, err := cond.GetConditions(server.DB, limit, offset)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	response.JSON(w, http.StatusOK, conds)
}
