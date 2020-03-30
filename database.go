package gorestframework

import (
	"log"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error
var IsDbInitialized = false

func InitDbConn(
	databaseDialect string,
	connectionString string,
	fnDbAutoMigrate func(db *gorm.DB),
) {
	log.Println("Opening db conn...")
	db, err = gorm.Open(databaseDialect, connectionString)
	if err != nil {
		log.Fatalln(err)
		panic("failed to connect database")
	}
	fnDbAutoMigrate(db)
}

func CloseDbConn() {
	db.Close()
	log.Println("Closing db conn...")
}

func DbOperation(fn func(db *gorm.DB)) {
	fn(db)
}
