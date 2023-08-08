package main

import (
	"os"

	"github.com/sirupsen/logrus"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/joho/godotenv"
	"github.com/kryptomind/bidboxapi/AccountsService/api/controllers"
	// BitgetWebSockets "github.com/kryptomind/bidboxapi/AccountsService/api/websockets"
)

var server = controllers.Server{}

func Run() {
	err := godotenv.Load()
	log := logrus.New()

	log.SetFormatter(&nested.Formatter{
		HideKeys:    true,
		FieldsOrder: []string{"file", "function"},
	})
	if err != nil {
		log.WithFields(logrus.Fields{
			"file":     "main.go",
			"function": "Run",
		}).Fatal("Error getting env")
	} else {
		log.WithFields(logrus.Fields{
			"file":     "main.go",
			"function": "Run",
		}).Info("Getting Values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	server.Run(":8080")
}

//	@title			BidBox Accounts Service API
//	@version		1.0
//	@description	This is BidBox Accounts Service.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/accounts

func main() {
	Run()
}
