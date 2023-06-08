package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GenerateBitgetSignature(apiSecret string, apiKey string, passphrase string, method string, uri string, timestamp string) string {

	message := fmt.Sprintf("%s%s%s", timestamp, method, uri)

	// Calculate HMAC-SHA256 signature
	hmac := hmac.New(sha256.New, []byte(apiSecret))
	hmac.Write([]byte(message))
	signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))

	return signature
}

func ValidateBitgetKeys(apiSecret string, apiKey string, passphrase string) (string, error) {
	host := "https://api.bitget.com"
	path := "/api/spot/v1/account/getInfo"
	url := host + path

	serverTimestamp := GetBitgetServerTimeStamp()

	signatures := GenerateBitgetSignature(apiSecret, apiKey, passphrase, "GET", path, serverTimestamp)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Print("-----error in request---", err.Error())

		return "", err
	}

	req.Header.Add("ACCESS-KEY", apiKey)
	req.Header.Add("ACCESS-PASSPHRASE", passphrase)
	req.Header.Add("ACCESS-TIMESTAMP", serverTimestamp)
	req.Header.Add("ACCESS-SIGN", signatures)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())

		return "", err
	}

	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print("---read error---", err.Error())
		return "", readErr
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("api key validation failed. Invalid credentials")

	}

	return string(body), nil
}

func GetAccountDetailsList(apiSecret string, apiKey string, passphrase string) (string, error) {
	host := "https://api.bitget.com"
	path := "/api/mix/v1/account/accounts?productType=sumcbl"
	url := host + path

	serverTimestamp := GetBitgetServerTimeStamp()
	signatures := GenerateBitgetSignature(apiSecret, apiKey, passphrase, "GET", path, serverTimestamp)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	req.Header.Add("ACCESS-KEY", apiKey)
	req.Header.Add("ACCESS-PASSPHRASE", passphrase)
	req.Header.Add("ACCESS-TIMESTAMP", serverTimestamp)
	req.Header.Add("ACCESS-SIGN", signatures)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())

		return "", err
	}

	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {

		return "", readErr
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", errors.New("unable to get account details at the moment")

	}

	return string(body), nil
}

func GetBitgetServerTimeStamp() string {

	url := "https://api.bitget.com/api/spot/v1/public/time"
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Print("-----error in request---", err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		fmt.Print("---read error---", err.Error())
	}

	var response bitgetServerTimeStampResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Print("---json unmarshal error---", err.Error())
	}

	return response.Data

}
