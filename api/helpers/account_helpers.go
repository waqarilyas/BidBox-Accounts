package helpers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/kryptomind/bidboxapi/AccountsService/api/models"
)

type Server struct {
	DB *gorm.DB
}

func (server *Server) UpdateUserBalanceByKeys(keyid uuid.UUID, apiSecret string, apiKey string, passphrase string) models.Accounts {
	decrypted_api_key, err := DecryptStrings(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	decrypted_passphrase, err := DecryptStrings(passphrase)
	if err != nil {
		log.Fatal(err)
	}
	decrypted_secret, err := DecryptStrings(apiSecret)
	if err != nil {
		log.Fatal(err)
	}

	response, err := GetAccountDetailsList(decrypted_secret, decrypted_api_key, decrypted_passphrase)
	if err != nil {
		fmt.Println("--something went wrong while getting user accoun details---", err)
	}

	accountResponse := AccountResponse{}
	responseBytes := []byte(response)

	error := json.Unmarshal(responseBytes, &accountResponse)
	if error != nil {
		fmt.Println(error)

	}

	accountData := accountResponse.Data[0]

	account := models.Accounts{
		CreatedAt:          time.Now().Format(time.RFC3339),
		UpdatedAt:          time.Now().Format(time.RFC3339),
		MarginCoin:         accountData.MarginCoin,
		AvailableBalance:   accountData.AvailableBalance,
		TotalMarginBalance: accountData.TotalMarginBalance,
		MarginValueUSDT:    accountData.MarginValueUSDT,
		MarginValueBTC:     accountData.MarginValueBTC,
		FloatingPnl:        accountData.FloatingPnl,
		ApiKeyId:           keyid,
	}

	if accountData.AvailableBalance == "" {
		account.AvailableBalance = "0"
	}

	if accountData.TotalMarginBalance == "" {
		account.TotalMarginBalance = "0"
	}
	if accountData.MarginValueUSDT == "" {
		account.MarginValueUSDT = "0"
	}
	if accountData.MarginValueBTC == "" {
		account.MarginValueBTC = "0"
	}
	if accountData.FloatingPnl == "" {
		account.FloatingPnl = "0"
	}

	accRes, err := account.SaveAccount(server.DB)
	if err != nil {
		fmt.Println(err)

	}

	// marshalledRes := json.Marshal(accRes)

	return accRes
}
