package db

import (
	"SpamhausResponseApi/model"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	mocket "github.com/selvatico/go-mocket"
)

func Connect() {
	var err error
	Conn, err = gorm.Open("sqlite3", "spamhausresponseapi.db")
	if err != nil {
		log.Fatalf("could not connected to database. %s", err)
	}

	Conn.Debug().AutoMigrate(model.IPDetails{})
}

func MockConnect() {
	var err error
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	mocket.Catcher.PanicOnEmptyResponse = true

	Conn, err = gorm.Open(mocket.DriverName, "connection_string")
	if err != nil {
		log.Fatalf("could not connected to database. %s", err)
	}
}

var Conn *gorm.DB
