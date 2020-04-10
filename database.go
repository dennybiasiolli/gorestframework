package gorestframework

import (
	"log"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error
var isDbInitialized = false

// IsDbInitialized returns if database is already initialized
func IsDbInitialized() bool {
	return isDbInitialized
}

// InitDbConn initialize a DB Connection using defined `databaseDialect` and `connectionString`.
// If defined, the `fnDbAutoMigrate` allows to execute auto-migration of database structure when the connection is opened.
// **Note**: don't forget to call `defer gorestframework.CloseDbConn()`
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
	isDbInitialized = true
	if fnDbAutoMigrate != nil {
		fnDbAutoMigrate(db)
	}
}

// CloseDbConn close database connection when is opened
func CloseDbConn() {
	if isDbInitialized {
		db.Close()
		log.Println("Closing db conn...")
	}
}

// DbOperation allows users to perform an operation on the active database.
func DbOperation(fn func(db *gorm.DB)) {
	if isDbInitialized {
		fn(db)
	}
}
