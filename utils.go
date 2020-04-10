package gorestframework

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// JSONContentType is a constant defining the content type for json requests
	JSONContentType string = "application/json"
)

// JSONRespond is a shortcut for JSONRespondWithStatus, using http.StatusOK as default status
func JSONRespond(w http.ResponseWriter, data interface{}) error {
	return JSONRespondWithStatus(w, data, http.StatusOK)
}

// JSONRespondWithStatus writes to a http.ResponseWriter a specific `data` value in JSON format, using the passed httpStatus
func JSONRespondWithStatus(w http.ResponseWriter, data interface{}, httpStatus int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	// data is a pointer to struct to marshal.
	// If no pointer is provided, an empty payload is returned.
	if data != nil {
		return json.NewEncoder(w).Encode(data)
	}
	return nil
}

// Respond answers providing the correct content based on client accept header.
func Respond(w http.ResponseWriter, r *http.Request, data interface{}) error {
	switch r.Header.Get("Accept") {
	case JSONContentType:
		return JSONRespond(w, data)
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		// it should be used a sentinel error instead: https://dave.cheney.net/tag/errors
		return fmt.Errorf("invalid Accept header")
	}
}

// RespondWithStatus answers providing the correct content based on client accept header,
// with custom HTTP status.
func RespondWithStatus(w http.ResponseWriter, r *http.Request, data interface{}, httpStatus int) error {
	switch r.Header.Get("Accept") {
	case JSONContentType:
		return JSONRespondWithStatus(w, data, httpStatus)
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return fmt.Errorf("invalid Accept header")
	}
}
