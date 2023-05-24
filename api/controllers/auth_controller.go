package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kryptomind/bidboxapi/AccountsService/api/auth"
	"github.com/kryptomind/bidboxapi/AccountsService/api/helpers"
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

	cipher, err := helpers.EncryptStrings(payload.Password)
	if err != nil {
		log.Fatal(err)
		return
	}

	newUser := admin.Admin{
		Email:    strings.ToLower(payload.Email),
		Password: cipher,
	}

	result := s.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		response.JSON(w, http.StatusConflict, "email already exist, please use another email address")
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
		response.JSON(w, http.StatusBadRequest, "Invalid email or Password")
		return
	}

	plain, err := helpers.DecryptStrings(user.Password)
	if err != nil {
		log.Fatal(err)
	}

	if plain != payload.Password {
		response.JSON(w, http.StatusBadRequest, "Invalid email or Password")
		return
	}

	userResponse := make(map[string]interface{})
	userResponse["id"] = user.Id.String()
	userResponse["email"] = user.Email
	userResponse["otp_verified"] = user.OtpVerified

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
		response.JSON(w, http.StatusBadRequest, "Invalid email or Password")
		return
	}

	dataToUpdate := admin.Admin{
		OtpSecret: key.Secret(),
		OtpUrl:    key.URL(),
	}

	s.DB.Model(&user).Updates(dataToUpdate)

	otpResponse := make(map[string]interface{})
	otpResponse["base32"] = key.Secret()
	otpResponse["otpauth_url"] = key.URL()
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
		response.JSON(w, http.StatusBadRequest, message)
		return
	}

	valid := totp.Validate(payload.Token, user.OtpSecret)
	if !valid {
		response.JSON(w, http.StatusBadRequest, message)
		return
	}

	dataToUpdate := admin.Admin{
		OtpVerified: true,
	}

	s.DB.Model(&user).Updates(dataToUpdate)

	token, err := auth.CreateToken(payload.UserId)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	userResponse := make(map[string]interface{})
	userResponse["id"] = user.Id.String()
	userResponse["email"] = user.Email
	userResponse["otp_verified"] = user.OtpVerified
	userResponse["token"] = token
	response.JSON(w, http.StatusOK, userResponse)
}

func (s *Server) ValidateOTP(w http.ResponseWriter, r *http.Request) {
	var payload *admin.OTPInput

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.JSON(w, http.StatusBadRequest, err)
		return
	}

	message := "Token is invalid or user doesn't exist"

	var user admin.Admin
	result := s.DB.First(&user, "id = ?", payload.UserId)
	if result.Error != nil {
		response.JSON(w, http.StatusBadRequest, message)
		return
	}

	valid := totp.Validate(payload.Token, user.OtpSecret)
	if !valid {
		response.JSON(w, http.StatusBadRequest, message)
		return
	}

	resp := make(map[string]bool)
	resp["otp_valid"] = true
	response.JSON(w, http.StatusOK, resp)
}
