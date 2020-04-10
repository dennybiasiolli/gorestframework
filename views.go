package gorestframework

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ViewInput contains the input of the View function
type ViewInput struct {
	Router     *mux.Router       // router used for adding the subroutes
	PathPrefix string            // prefix where the subroutes will be added
	Controller *ControllerOutput // optional custom controller for the routes
	ModelPtr   interface{}       // model used for reading/writing data from/to database
}

// ViewOutput contains the output of the View function
type ViewOutput struct {
	Router     *mux.Router      // router used for adding the subroutes
	Subrouter  *mux.Router      // subrouter containing the subroutes
	PathPrefix string           // prefix where the subroutes are added
	Controller ControllerOutput // controller used for the routes
	ModelPtr   interface{}      // model used for reading/writing data from/to database
}

// View returns a ViewOutput containing all routes for a specific PathPrefix
func View(
	input *ViewInput,
) ViewOutput {
	var controller ControllerOutput
	if input.Controller == nil {
		controller = Controller(input.ModelPtr)
	} else {
		controller = *input.Controller
	}
	if strings.HasPrefix(input.PathPrefix, "/") == false {
		input.PathPrefix = fmt.Sprintf("/%s", input.PathPrefix)
	}

	subrouter := input.Router.PathPrefix(input.PathPrefix).Subrouter()
	subrouter.
		HandleFunc("/", controller.GetAll).
		Methods(http.MethodGet)

	subrouter.
		HandleFunc("/{id:[0-9]+}/", controller.Get).
		Methods(http.MethodGet)

	subrouter.
		HandleFunc("/", controller.Post).
		Methods(http.MethodPost)

	subrouter.
		HandleFunc("/{id:[0-9]+}/", controller.Put).
		Methods(http.MethodPut)
	subrouter.
		HandleFunc("/{id:[0-9]+}/", controller.Patch).
		Methods(http.MethodPatch)

	subrouter.
		HandleFunc("/{id:[0-9]+}/", controller.Delete).
		Methods(http.MethodDelete)

	return ViewOutput{
		Router:     input.Router,
		Subrouter:  subrouter,
		PathPrefix: input.PathPrefix,
		Controller: controller,
		ModelPtr:   input.ModelPtr,
	}
}
