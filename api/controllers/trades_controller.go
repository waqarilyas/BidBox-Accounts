package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"strconv"

	helpers "github.com/ahmed-023/bitget-helpers"
)

type Positions struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int    `json:"requestTime"`
	Data        []struct {
		MarginCoin        string `json:"marginCoin"`
		Symbol            string `json:"symbol"`
		HoldSide          string `json:"holdSide"`
		OpenDelegateCount string `json:"openDelegateCount"`
		Margin            string `json:"margin"`
		Available         string `json:"available"`
		Locked            string `json:"locked"`
		Total             string `json:"total"`
		Leverage          int    `json:"leverage"`
		AchievedProfits   string `json:"achievedProfits"`
		AverageOpenPrice  string `json:"averageOpenPrice"`
		MarginMode        string `json:"marginMode"`
		HoldMode          string `json:"holdMode"`
		UnrealizedPL      string `json:"unrealizedPL"`
		LiquidationPrice  string `json:"liquidationPrice"`
		KeepMarginRate    string `json:"keepMarginRate"`
		MarketPrice       string `json:"marketPrice"`
		CTime             string `json:"cTime"`
	} `json:"data"`
}

func (server *Server) GetActiveTrades(w http.ResponseWriter, r *http.Request) {
	// Make API call to get all positions
	// ...
	api_key := os.Getenv("API_KEY")
	secret_key := os.Getenv("SECRET_KEY")
	passphrase := os.Getenv("PASSPHRASE")

	host := "https://api.bitget.com"
	path := "/api/mix/v1/position/allPosition?productType=sumcbl"
	url := host + path
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	server_time := helpers.GetBitgetServerTimeStamp()
	signatures := helpers.GenerateBitgetSignature(secret_key, api_key, passphrase, "GET", path, server_time)

	req.Header.Add("ACCESS-KEY", api_key)
	req.Header.Add("ACCESS-SIGN", signatures)
	req.Header.Add("ACCESS-TIMESTAMP", server_time)
	req.Header.Add("ACCESS-PASSPHRASE", passphrase)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	// Parse the API response into a AllPositionsResponse struct
	var resp Positions
	if err := json.Unmarshal(body, &resp); err != nil {
		// handle error
		log.Fatal("parse error")
	}

	// Iterate through each position in the response and print whether it's open or closed
	for _, pos := range resp.Data {
		if pos.HoldSide == "long" || pos.HoldSide == "short" {
			if total, err := strconv.ParseFloat(pos.Total, 64); err == nil && total > 0 {
				fmt.Printf("%s %s position is open\n", pos.HoldSide, pos.Symbol)
			} else {
				fmt.Printf("%s %s position is closed\n", pos.HoldSide, pos.Symbol)
			}
		}
	}

	// ...
}
