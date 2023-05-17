package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/kryptomind/bidboxapi/AccountsService/api/models/admin"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
	"github.com/pquerna/otp/totp"
)

func (s *Server) SignUpUser(w http.ResponseWriter, r *http.Request) {
	var payload *admin.RegisterUserInput

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.JSON(w, http.StatusBadRequest, err)
		return
	}

	newUser := admin.Admin{
		Email:    strings.ToLower(payload.Email),
		Password: payload.Password,
	}

	result := s.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		response.JSON(w, http.StatusConflict, errors.New("Email already exist, please use another email address"))
		return
	} else if result.Error != nil {
		response.JSON(w, http.StatusBadGateway, result.Error)
		return
	}

	response.JSON(w, http.StatusCreated, "Registered successfully, please login")
}

func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	var payload *admin.LoginUserInput

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.JSON(w, http.StatusBadRequest, err)
		return
	}

	var user admin.Admin
	result := s.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		response.JSON(w, http.StatusBadRequest, errors.New("Invalid email or Password"))
		return
	}

	userResponse := map[string]interface{}{
		"id":           user.Id.String(),
		"email":        user.Email,
		"otp_verified": user.OtpVerified,
	}
	response.JSON(w, http.StatusOK, userResponse)
}

func (s *Server) GenerateOTP(w http.ResponseWriter, r *http.Request) {
	var payload *admin.OTPInput

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.JSON(w, http.StatusBadRequest, err)
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "codevoweb.com",
		AccountName: "admin@admin.com",
		SecretSize:  15,
	})

	if err != nil {
		panic(err)
	}

	var user admin.Admin
	result := s.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		response.JSON(w, http.StatusBadRequest, errors.New("Invalid email or Password"))
		return
	}

	dataToUpdate := admin.Admin{
		Otp_secret:   key.Secret(),
		Otp_auth_url: key.URL(),
	}

	s.DB.Model(&user).Updates(dataToUpdate)

	otpResponse := map[string]interface{}{
		"base32":      key.Secret(),
		"otpauth_url": key.URL(),
	}
	response.JSON(w, http.StatusOK, otpResponse)
}

func (s *Server) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var payload *admin.OTPInput

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.JSON(w, http.StatusBadRequest, err)
		return
	}

	message := "Token is invalid or user doesn't exist"

	var user admin.Admin
	result := s.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		response.JSON(w, http.StatusBadRequest, errors.New(message))
		return
	}

	valid := totp.Validate(payload.Token, user.Otp_secret)
	if !valid {
		response.JSON(w, http.StatusBadRequest, errors.New(message))
		return
	}

	dataToUpdate := admin.Admin{
		OtpVerified: true,
	}

	s.DB.Model(&user).Updates(dataToUpdate)

	userResponse := map[string]interface{}{
		"id":           user.Id.String(),
		"email":        user.Email,
		"otp_verified": user.OtpVerified,
	}
	response.JSON(w, http.StatusOK, userResponse)
}

func (s *Server) ValidateOTP(w http.ResponseWriter, r *http.Request) {
	var payload *admin.OTPInput

	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		response.JSON(w, http.StatusBadRequest, err)
		return
	}

	message := "Token is invalid or user doesn't exist"

	var user admin.Admin
	result := s.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		response.JSON(w, http.StatusBadRequest, errors.New(message))
		return
	}

	valid := totp.Validate(payload.Token, user.Otp_secret)
	if !valid {
		response.JSON(w, http.StatusBadRequest, errors.New(message))
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{"otp_valid": true})
}
