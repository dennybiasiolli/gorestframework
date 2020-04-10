package gorestframework

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type ViewInput struct {
	Router     *mux.Router
	PathPrefix string
	Controller *ControllerOutput
	ModelPtr   interface{}
}

type ViewOutput struct {
	Router     *mux.Router
	Subrouter  *mux.Router
	PathPrefix string
	Controller ControllerOutput
	ModelPtr   interface{}
}

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
