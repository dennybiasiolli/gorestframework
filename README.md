# Go Rest Framework

[![Go Report Card](https://goreportcard.com/badge/github.com/dennybiasiolli/gorestframework)](https://goreportcard.com/report/github.com/dennybiasiolli/gorestframework)

#### Install module

`go get -u github.com/dennybiasiolli/gorestframework`


#### Usage example

```go
package main

import (
	"github.com/dennybiasiolli/gorestframework"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// create the model definition
// more info here: https://gorm.io/docs/models.html
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// create migration function
func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(
		// passing all desired models here
		&Product{},
	)
}

// create SetView function
func SetViews(router *mux.Router) {
	gorestframework.View(&gorestframework.ViewInput{
		Router:     router,
		PathPrefix: "/products",
		ModelPtr:   &Product{},
	})
}

func main() {
	// initializing database connection
	gorestframework.InitDbConn(
		"sqlite3",  // DatabaseDialect,
		"test.db",  // DatabaseConnectionString,
		MigrateModels,
	)
	defer gorestframework.CloseDbConn()

	// start HTTP listener
	gorestframework.StartHTTPListener(
		true,  // RouterActivateLog,
		true,  // RouterUseCORS,
		views.SetViews,
	)
}
```
