package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/go-playground/validator/v10"
	"github.com/kryptomind/bidboxapi/AccountsService/api/auth"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

type emailNext func(http.ResponseWriter, *http.Request, string)

var validate = validator.New()

func MiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		next(w, r)
	}
}

func ValidateEmail(next emailNext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set("Content-Type", "application/json")

		email := r.URL.Query().Get("email")
		if email == "" {
			response.ERROR(w, http.StatusBadRequest, errors.New("email is required"))
			return
		}

		if !govalidator.IsEmail(email) {
			response.ERROR(w, http.StatusBadRequest, errors.New("invalid email address"))
			return
		}
		next(w, r, email)
	}
}

func ValidateBody(next http.HandlerFunc, v interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := validate.Struct(v); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if r.ContentLength == 0 {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		if err := json.Unmarshal(body, &v); err != nil {
			fmt.Println(err)
			http.Error(w, "Error unmarshaling request body", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(v); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		r.Header.Set("Content-Type", "application/json")

		next(w, r)
		// next.ServeHTTP(w, r)
	}
}

func MiddlewareAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			response.ERROR(w, http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}
		next(w, r)
	}
}
