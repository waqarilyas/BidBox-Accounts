package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kryptomind/bidboxapi/AccountsService/api/auth"
	"github.com/kryptomind/bidboxapi/AccountsService/api/helpers"
	"github.com/kryptomind/bidboxapi/AccountsService/api/models/admin"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
	"github.com/pquerna/otp/totp"
	// "github.com/dgrijalva/jwt-go"
)

type UserBody struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"123"`
}

var (
	user = "ahmed023"
	pass = "monke"
)

// Generate JWT Token godoc
// @Summary      Generate JWT Token
// @Description  generate token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  UserBody  true  "user body"
// @Success      200  {object}  admin.CoinPair
// @Failure      422  {string}  coin required
// @Failure      500  {string}  server error
// @Router       /generate-token [post]
func (s *Server) GenerateJWT(w http.ResponseWriter, r *http.Request) {
	userbody := UserBody{}
	if err := json.NewDecoder(r.Body).Decode(&userbody); err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	if userbody.Username != user || userbody.Password != pass {
		response.ERROR(w, http.StatusBadRequest, errors.New("username or password incorrect"))
		return
	}

	token, err := auth.CreateToken(userbody.Username + userbody.Password)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res := make(map[string]interface{}, 0)
	res["token"] = token

	response.JSON(w, http.StatusOK, res)
}

// Sign Up New Admin User godoc
// @Summary      Sign Up Admin User
// @Description  sign up admin user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  admin.RegisterUserInput true  "user body"
// @Success      201  {string}  Registered succesfully
// @Failure      400  {string}  coin required
// @Failure      409  {string}  coin required
// @Failure      502  {string}  server error
// @Router       /admin/register [post]
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

type LoginResp struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	OtpVerified bool   `json:"otp_verified"`
	JwtToken    string `json:"jwt_token"`
}

// Login Admin godoc
// @Summary      Login Admin
// @Description  login admin
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  admin.LoginUserInput true  "user login body"
// @Success      200  {object} LoginResp
// @Failure      400  {string}  coin required
// @Failure      500  {string}  server error
// @Router       /admin/login [post]
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

	fmt.Println(plain)
	if plain != payload.Password {
		response.JSON(w, http.StatusBadRequest, "Invalid email or Password")
		return
	}
	// fmt.Println("user", user.Id.String())
	// logic to create secret key
	token, err := auth.CreateToken(user.Id.String())
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	userResponse := make(map[string]interface{})
	userResponse["id"] = user.Id.String()
	userResponse["email"] = user.Email
	userResponse["otp_verified"] = user.OtpVerified
	userResponse["jwt_token"] = token
	response.JSON(w, http.StatusOK, userResponse)
}

type OtpGenResp struct {
	Base32     string `json:"base32"`
	OtpauthUrl string `json:"otpauth_url"`
}

// Generate OTP godoc
// @Summary      Generate OTP
// @Description  generate otp
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  admin.OTPInput true  "otp input"
// @Success      200  {object} OtpGenResp
// @Failure      400  {string}  coin required
// @Failure      500  {string}  server error
// @Router       /admin/otp/generate [post]
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

type OtpVerifyResp struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	OtpVerified bool   `json:"otp_verified"`
	Token       string `json:"token"`
}

// Verify OTP godoc
// @Summary      Verify OTP
// @Description  verify otp
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  admin.OTPInput true  "otp input"
// @Success      200  {object} OtpVerifyResp
// @Failure      400  {string}  coin required
// @Failure      500  {string}  server error
// @Router       /admin/otp/verify [post]
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

type OtpValidResp struct {
	OtpValid bool `json:"otp_valid"`
}

// Validate OTP godoc
// @Summary      Validate OTP
// @Description  validate otp
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  admin.OTPInput true  "otp input"
// @Success      200  {object} OtpValidResp
// @Failure      400  {string}  coin required
// @Router       /admin/otp/validate [post]
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

// Enable OTP godoc
// @Summary      Enable OTP
// @Description  enable otp
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        otp  body  admin.OTPResponse true  "otp enable"
// @Param        id  query  string true  "user id"
// @Success      200  {object} admin.Resp
// @Failure      400  {string}  coin required
// @Failure      500  {string}  coin required
// @Router       /admin/otp/enable [post]
func (s *Server) EnableOTP(w http.ResponseWriter, r *http.Request) {
	var payload *admin.OTPResponse

	id := r.URL.Query().Get("id")
	if id == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("id is required as query param"))
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if payload.OPTEnabled == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("otp enabled required"))
		return
	}

	var user admin.Admin

	result := s.DB.First(&user, "id = ?", strings.ToLower(id))
	if result.Error != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	enable := payload.OPTEnabled
	dataToUpdate := admin.Admin{
		OtpEnabled: enable,
	}
	uotp, err := dataToUpdate.UpdateOtp(s.DB, enable)

	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	response.JSON(w, http.StatusOK, uotp)

}

// Change Password godoc
// @Summary      Change Password
// @Description  change password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        otp  body  admin.ChangePasswordInput true  "change password"
// @Param        id  query  string true  "user id"
// @Param        Authorization  header  string  true  "Authorization"
// @Success      200  {string} kksks
// @Failure      400  {string}  coin required
// @Router       /admin/changePassword [post]
func (s *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("id is required as query param"))
		return
	}
	var payload *admin.ChangePasswordInput

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.JSON(w, http.StatusBadRequest, err)
		return
	}

	if payload.NewPassword == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("new password is required"))
		return
	}

	if payload.Password == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("password is required"))
		return
	}

	// tokenID, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	// if tokenID != id {
	// 	response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }

	var user admin.Admin
	result := s.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		response.JSON(w, http.StatusBadRequest, "Invalid Password")
		return
	}

	plain, err := helpers.DecryptStrings(user.Password)
	if err != nil {
		log.Fatal(err)
	}

	if plain != payload.Password {
		response.JSON(w, http.StatusBadRequest, "Invalid Password")
		return
	}
	cipher, err := helpers.EncryptStrings(payload.NewPassword)
	if err != nil {
		log.Fatal(err)
		return
	}
	dataToUpdate := admin.Admin{
		Password: cipher,
	}

	s.DB.Model(&user).Updates(dataToUpdate)

	response.JSON(w, http.StatusOK, "Password successfully updated")
}

// Change Timeframe godoc
// @Summary      Change Timeframe
// @Description  change timeframe
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        otp  body  admin.ChangeTimeframeInput true  "change timeframe"
// @Param        id  query  string true  "user id"
// @Param        Authorization  header  string  true  "Authorization"
// @Success      200  {string} kksks
// @Failure      400  {string}  coin required
// @Router       /admin/changetimeframe [post]
func (s *Server) changeTimeframe(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")
	if id == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("id is required as query param"))
		return
	}
	var payload *admin.ChangeTimeframeInput

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		response.JSON(w, http.StatusBadRequest, err)
		return
	}

	if payload.Timeframe == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("timeframe is required"))
		return
	}
	// tokenID, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	// if tokenID != id {
	// 	response.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }
	var user admin.Settings
	// result := s.DB.First(&user, "id = ?", id)
	// if result.Error != nil {
	// 	response.JSON(w, http.StatusBadRequest, "Invalid Credentials")
	// 	return
	// }

	dataToUpdate := admin.Settings{
		Timeframe: payload.Timeframe,
	}

	s.DB.Model(&user).Updates(dataToUpdate)

	response.JSON(w, http.StatusOK, "TimeFrame successfully updated")
}

type TimeframeResp struct {
	Timeframe string `json:"timeframe"`
}

// Get Timeframe godoc
// @Summary      Get Timeframe
// @Description  get timeframe
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        id  query  string true  "user id"
// @Param        Authorization  header  string  true  "Authorization"
// @Success      200  {object} TimeframeResp
// @Failure      400  {string}  coin required
// @Failure      500  {string}  coin required
// @Router       /admin/timeframe [post]
func (server *Server) GetTimeframe(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("id is required as query param"))
		return
	}
	key := &admin.Settings{}

	timeframe, err := key.GetTimeframe(server.DB)
	if err != nil {
		response.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	res := make(map[string]interface{})
	res["timeframe"] = timeframe
	response.JSON(w, http.StatusOK, res)
}
