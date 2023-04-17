package controllers

import (
	"net/http"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "")
}


func (server *Server) GetUserAcounts(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, "")
}


