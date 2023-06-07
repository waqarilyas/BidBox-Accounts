package helpers

type bitgetServerTimeStampResponse struct {
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	RequestTime int    `json:"requestTime"`
	Data        string `json:"data"`
}

type AccountResponse struct {
	Code        string        `json:"code"`
	Message     string        `json:"msg"`
	RequestTime int64         `json:"requestTime"`
	Data        []AccountData `json:"data"`
}

type AccountData struct {
	MarginCoin         string `json:"marginCoin"`
	AvailableBalance   string `json:"available"`
	TotalMarginBalance string `json:"equity"`
	MarginValueUSDT    string `json:"usdtEquity"`
	MarginValueBTC     string `json:"btcEquity"`
	FloatingPnl        string `json:"unrealizedPL"`
}
