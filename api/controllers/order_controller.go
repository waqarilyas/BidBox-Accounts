package controllers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	helpers "github.com/ahmed-023/bitget-helpers"
	"github.com/asaskevich/govalidator"
	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
	"github.com/kryptomind/bidboxapi/AccountsService/api/response"
)

type Order struct {
	Symbol           string `json:"symbol"`
	MarginCoin       string `json:"marginCoin"`
	Size             string `json:"size"`
	Side             string `json:"side"`
	OrderType        string `json:"orderType"`
	TimeInForceValue string `json:"timeInForceValue"`
}

type OrderResponse struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int64  `json:"requestTime"`
	Data        struct {
		ClientOid string `json:"clientOid"`
		OrderID   string `json:"orderId"`
	} `json:"data"`
}

func (server *Server) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	log.Print(order.Side)

	email := r.URL.Query().Get("email")
	if email == "" {
		response.ERROR(w, http.StatusBadRequest, errors.New("email is required"))
		return
	}

	if !govalidator.IsEmail(email) {
		response.ERROR(w, http.StatusBadRequest, errors.New("invalid email address"))
		return
	}

	//get keys by user id
	key := models.Key{}
	keys, err := key.FindKeysByEmail(server.DB, email)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, errors.New("User not found"))
		return
	}

	var api_key string
	var secret_key string
	var passphrase string
	for _, v := range *keys {
		if strings.ToLower(v.Service) != "bitget" {
			response.JSON(w, http.StatusNoContent, errors.New("Exchange coming soon"))
			return
		} else {
			api_key, err = helpers.DecryptStrings(v.ApiKey)
			if err != nil {
				log.Fatal(err)
				return
			}
			secret_key, err = helpers.DecryptStrings(v.SecretKey)
			if err != nil {
				log.Fatal(err)
				return
			}
			passphrase, err = helpers.DecryptStrings(v.Passphrase)
			if err != nil {
				log.Fatal(err)
				return
			}
		}
	}
	if api_key == "" {
		response.ERROR(w, http.StatusNoContent, errors.New("api key not found"))
		return
	}
	var orderResp OrderResponse
	str := NewOrder(api_key, secret_key, passphrase, &order)
	err = json.Unmarshal([]byte(str), &orderResp)
	if err != nil {
		response.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Write([]byte(orderResp.Data.OrderID))
}

func GenerateBitgetSignature(apiSecret string, apiKey string, passphrase string, method string, uri string, timestamp string, requestBody string) string {

	message := ""
	if method == "GET" {
		message = fmt.Sprintf("%s%s%s", timestamp, method, uri)
	} else if method == "POST" {
		message = fmt.Sprintf("%s%s%s%s", timestamp, method, uri, requestBody)
	}

	// Calculate HMAC-SHA256 signature
	hmac := hmac.New(sha256.New, []byte(apiSecret))
	hmac.Write([]byte(message))
	signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))

	return signature
}

func NewOrder(api_key string, secret_key string, passphrase string, order *Order) string {

	host := "https://api.bitget.com"
	path := "/api/mix/v1/order/placeOrder"
	url := host + path

	method := "POST"
	client := &http.Client{}

	jsonVal, _ := json.Marshal(order)

	server_time := helpers.GetBitgetServerTimeStamp()
	signatures := GenerateBitgetSignature(secret_key, api_key, passphrase, "POST", path, server_time, string(jsonVal))

	log.Println(server_time)
	log.Println(signatures)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonVal))
	req.Header.Add("ACCESS-KEY", api_key)
	req.Header.Add("ACCESS-SIGN", signatures)
	req.Header.Add("ACCESS-TIMESTAMP", server_time)
	req.Header.Add("ACCESS-PASSPHRASE", passphrase)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("local", "zh-CN")

	if err != nil {
		log.Fatal(err)
		return ""
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(body)
}
