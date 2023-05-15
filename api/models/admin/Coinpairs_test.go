package admin

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

type TestDB struct {
	DB *gorm.DB
}

func initDB() *gorm.DB {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Fatal(err)
	}

	Dbdriver := os.Getenv("DB_DRIVER")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbPort := os.Getenv("DB_PORT")
	DbHost := os.Getenv("DB_HOST")
	DbName := os.Getenv("DB_NAME")

	var err error
	db := TestDB{}
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		//fmt.Sprintln("Cannot connect to the %s database", Dbdriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Sprintln("Connected to the database")
	}

	return db.DB
}

func TestGetCoins(t *testing.T) {

	db := initDB()
	coinpairs := CoinPair{}

	c, err := coinpairs.GetAllCoins(db)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range *c {
		log.Println(v.Coin)
	}
}
