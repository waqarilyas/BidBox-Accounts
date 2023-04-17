package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"strconv"

	helpers "github.com/ahmed-023/bitget-helpers"
)

type Data struct {
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
}

type Positions struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int    `json:"requestTime"`
	Data        []Data `json:"data"`
}

func getTradeData() ([]Data, error) {
	// Make API call to get all positions
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
		return []Data{}, err
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
		return []Data{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []Data{}, err
	}

	// Parse the API response into a AllPositionsResponse struct
	var resp Positions
	if err := json.Unmarshal(body, &resp); err != nil {
		// handle error
		return []Data{}, err
	}

	return resp.Data, nil
}

func (server *Server) GetOpenTrades(w http.ResponseWriter, r *http.Request) {
	resp, err := getTradeData()

	if err != nil {
		log.Fatal(err)
	}

	var open_trades []Data
	// Iterate through each position in the response and print whether it's open or closed
	for _, pos := range resp {
		if pos.HoldSide == "long" || pos.HoldSide == "short" {
			if total, err := strconv.ParseFloat(pos.Total, 64); err == nil && total > 0 {
				open_trades = append(open_trades, pos)
			}
		}
	}

	json, err := json.Marshal(open_trades)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (server *Server) GetClosedTrades(w http.ResponseWriter, r *http.Request) {
	resp, err := getTradeData()

	if err != nil {
		log.Fatal(err)
	}

	var closed_trades []Data
	// Iterate through each position in the response and print whether it's open or closed
	for _, pos := range resp {
		if pos.HoldSide == "long" || pos.HoldSide == "short" {
			if total, err := strconv.ParseFloat(pos.Total, 64); err == nil && total > 0 {
				continue
			} else {
				closed_trades = append(closed_trades, pos)
			}
		}
	}

	json, err := json.Marshal(closed_trades)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
