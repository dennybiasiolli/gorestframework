/*
Package gorestframework implements a simple library for creating REST endpoints in an easy way.

An example could be the following:

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
*/
package gorestframework

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// StartHTTPListener starts the HTTP Listener.
// HOST and PORT can be passed via ENV, unless the default are
// - HOST=localhost
// - PORT=8000
//
// activateLog it's used to log requests to console
// useCORS enable CORS capability for all hosts, PR are welcome!
// fnSetViews is a callback function for configuring the mux.Router
func StartHTTPListener(
	activateLog bool,
	useCORS bool,
	fnSetViews func(*mux.Router),
) {
	var wait time.Duration
	flag.DurationVar(
		&wait,
		"graceful-timeout",
		time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m",
	)
	flag.Parse()

	router := mux.NewRouter()
	if activateLog {
		router.Use(LoggingMiddleware)
	}
	if useCORS {
		router.Methods(http.MethodOptions)
		router.Use(CORSMiddleware)
	}

	fnSetViews(router)

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%s", host, port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// run our server in a goroutine so that it doesn't block.
	go func() {
		log.Println("Starting http.Server on", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("Shutting down")

	// os.Exit(0)
}
