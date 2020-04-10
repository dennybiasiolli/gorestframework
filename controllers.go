package gorestframework

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// ControllerOutput contains the output of the Controller function
type ControllerOutput struct {
	ModelPtr interface{}                                  // pointer to the defined model
	GetAll   func(w http.ResponseWriter, r *http.Request) // function for retreiving all records
	Get      func(w http.ResponseWriter, r *http.Request) // function for retreiving a single record
	Post     func(w http.ResponseWriter, r *http.Request) // function for adding a record
	Put      func(w http.ResponseWriter, r *http.Request) // function for updating a record
	Patch    func(w http.ResponseWriter, r *http.Request) // function for updating a record
	Delete   func(w http.ResponseWriter, r *http.Request) // function for deleting a record
}

// Controller returns a ControllerOutput containing all REST handlers
//
// controllerName contains
func Controller(modelPtr interface{}) ControllerOutput {
	T := reflect.TypeOf(modelPtr)
	if T.Kind() == reflect.Ptr {
		T = T.Elem()
	}
	return ControllerOutput{
		ModelPtr: modelPtr,

		GetAll: func(w http.ResponseWriter, r *http.Request) {
			DbOperation(func(db *gorm.DB) {
				resultsPtr := reflect.New(reflect.SliceOf(T)).Interface()
				db.Find(resultsPtr)
				JsonRespond(w, resultsPtr)
			})
		},

		Get: func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			DbOperation(func(db *gorm.DB) {
				resultPtr := reflect.New(T).Interface()
				res := db.First(resultPtr, vars["id"])
				if res.RecordNotFound() {
					JsonRespondWithStatus(w, map[string]interface{}{
						"error":   true,
						"message": fmt.Sprintf("%s with ID %s not found", T.Name(), vars["id"]),
					}, http.StatusNotFound)
					return
				} else if err := res.Error; err != nil {
					JsonRespondWithStatus(w, err, http.StatusBadRequest)
					return
				}
				JsonRespond(w, resultPtr)
			})
		},

		Post: func(w http.ResponseWriter, r *http.Request) {
			// Declare a new struct.
			modelDataPtr := reflect.New(T).Interface()

			// Try to decode the request body into the struct. If there is an error,
			// respond to the client with the error message and a 400 status code.
			err := json.NewDecoder(r.Body).Decode(&modelDataPtr)
			if err != nil {
				JsonRespondWithStatus(w, map[string]interface{}{
					"Severity": "error",
					"Message":  "Unable to parse JSON body.",
				}, http.StatusBadRequest)
				return
			}

			// try to create the record into the database
			DbOperation(func(db *gorm.DB) {
				if err := db.Create(modelDataPtr).Error; err != nil {
					JsonRespondWithStatus(w, err, http.StatusBadRequest)
					return
				}
				JsonRespond(w, modelDataPtr)
			})
		},

		Put: func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)

			// Declare a new struct.
			modelDataPtr := reflect.New(T).Interface()
			newModelDataPtr := reflect.New(T).Interface()

			// Try to decode the request body into the struct. If there is an error,
			// respond to the client with the error message and a 400 status code.
			err := json.NewDecoder(r.Body).Decode(&newModelDataPtr)
			if err != nil {
				JsonRespondWithStatus(w, map[string]interface{}{
					"Severity": "error",
					"Message":  "Unable to parse JSON body.",
				}, http.StatusBadRequest)
				return
			}

			// try to update the record into the database
			DbOperation(func(db *gorm.DB) {
				res := db.First(modelDataPtr, vars["id"])
				if res.RecordNotFound() {
					JsonRespondWithStatus(w, map[string]interface{}{
						"error":   true,
						"message": fmt.Sprintf("%s with ID %s not found", T.Name(), vars["id"]),
					}, http.StatusNotFound)
					return
				} else if err := res.Error; err != nil {
					JsonRespondWithStatus(w, err, http.StatusBadRequest)
					return
				}

				id := reflect.ValueOf(modelDataPtr).Elem().FieldByName("ID").Int()
				newIDField := reflect.ValueOf(newModelDataPtr).Elem().FieldByName("ID")
				if newIDField.IsValid() && id != newIDField.Int() {
					JsonRespondWithStatus(w, map[string]interface{}{
						"error": true,
						"message": fmt.Sprintf(
							"Unable to change ID of %s from %v to %v. PRs are welcome!",
							T.Name(), id, newIDField,
						),
					}, http.StatusBadRequest)
					return
				}

				res = db.Model(modelDataPtr).Updates(newModelDataPtr)
				if err := res.Error; err != nil {
					JsonRespondWithStatus(w, err, http.StatusBadRequest)
					return
				}
				JsonRespond(w, modelDataPtr)
			})
		},

		Patch: func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)

			// Declare a new struct.
			modelDataPtr := reflect.New(T).Interface()
			newModelDataPtr := reflect.New(T).Interface()

			// Try to decode the request body into the struct. If there is an error,
			// respond to the client with the error message and a 400 status code.
			err := json.NewDecoder(r.Body).Decode(&newModelDataPtr)
			if err != nil {
				JsonRespondWithStatus(w, map[string]interface{}{
					"Severity": "error",
					"Message":  "Unable to parse JSON body.",
				}, http.StatusBadRequest)
				return
			}

			// try to update the record into the database
			DbOperation(func(db *gorm.DB) {
				res := db.First(modelDataPtr, vars["id"])
				if res.RecordNotFound() {
					JsonRespondWithStatus(w, map[string]interface{}{
						"error":   true,
						"message": fmt.Sprintf("%s with ID %s not found", T.Name(), vars["id"]),
					}, http.StatusNotFound)
					return
				} else if err := res.Error; err != nil {
					JsonRespondWithStatus(w, err, http.StatusBadRequest)
					return
				}

				id := reflect.ValueOf(modelDataPtr).Elem().FieldByName("ID").Int()
				newIDField := reflect.ValueOf(newModelDataPtr).Elem().FieldByName("ID")
				if newIDField.IsValid() && id != newIDField.Int() {
					JsonRespondWithStatus(w, map[string]interface{}{
						"error": true,
						"message": fmt.Sprintf(
							"Unable to change ID of %s from %v to %v. PRs are welcome!",
							T.Name(), id, newIDField,
						),
					}, http.StatusBadRequest)
					return
				}

				res = db.Model(modelDataPtr).Updates(newModelDataPtr)
				if err := res.Error; err != nil {
					JsonRespondWithStatus(w, err, http.StatusBadRequest)
					return
				}
				JsonRespond(w, modelDataPtr)
			})
		},

		Delete: func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)

			// Declare a new struct.
			modelDataPtr := reflect.New(T).Interface()

			// try to delete the record into the database
			DbOperation(func(db *gorm.DB) {
				res := db.First(modelDataPtr, vars["id"])
				if res.RecordNotFound() {
					JsonRespondWithStatus(w, map[string]interface{}{
						"error":   true,
						"message": fmt.Sprintf("%s with ID %s not found", T.Name(), vars["id"]),
					}, http.StatusNotFound)
					return
				} else if err := res.Error; err != nil {
					JsonRespondWithStatus(w, err, http.StatusBadRequest)
					return
				}

				res = db.Model(modelDataPtr).Delete(modelDataPtr)
				if err := res.Error; err != nil {
					JsonRespondWithStatus(w, err, http.StatusBadRequest)
					return
				}
				JsonRespond(w, modelDataPtr)
			})
		},
	}
}
