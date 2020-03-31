# Go Rest Framework

Readme coming soon...

#### Example

`views/views.go`

```go
package views

import (
	"go-rest-service/models"

	"github.com/dennybiasiolli/gorestframework"
	"github.com/gorilla/mux"
)

func SetViews(router *mux.Router) {
	gorestframework.View(&gorestframework.ViewInput{
		Router:     router,
		PathPrefix: "/products",
		ModelPtr:   &models.Product{},
	})
}
```

`models/product.go`

```go
package models

import "github.com/jinzhu/gorm"

type Product struct {
	gorm.Model
	Code  string
	Price uint
}
```

`models/models.go`

```go
package models

import "github.com/jinzhu/gorm"

func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(
		&Product{},
	)
}
```

`main.go`

```go
package main

import (
	"go-rest-service/models"
	"go-rest-service/views"

	"github.com/dennybiasiolli/gorestframework"
)

func main() {
	gorestframework.InitDbConn(
		"sqlite3",  // DatabaseDialect,
		"test.db",  // DatabaseConnectionString,
		models.MigrateModels,
	)
	defer gorestframework.CloseDbConn()

	gorestframework.StartHTTPListener(
		true,  // RouterActivateLog,
		true,  // RouterUseCORS,
		views.SetViews,
	)
}

```
