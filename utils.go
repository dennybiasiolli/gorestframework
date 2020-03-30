package gorestframework

import (
	"encoding/json"
	"net/http"
)

func JsonRespond(w http.ResponseWriter, data interface{}) {
	JsonRespondWithStatus(w, data, http.StatusOK)
}
func JsonRespondWithStatus(w http.ResponseWriter, data interface{}, httpStatus int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)
}
