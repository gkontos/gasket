package aceweb

import (
	"encoding/json"
	"net/http"

	log "github.com/gkontos/gasket/acelog"
	service "github.com/gkontos/gasket/aceservice"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

type ValidationError struct {
	Err     error
	Message string
}

func (e *ValidationError) Error() string {
	if e.Message != "" {
		return e.Message + " " + e.Err.Error()
	}
	return e.Err.Error()
}

// ReturnErrorJSON will create and return a json encoded exception
func ReturnErrorJSON(w http.ResponseWriter, err error) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	httpStatus := http.StatusInternalServerError
	if err != nil {
		switch err.(type) {
		case *RequestParseError:
			httpStatus = http.StatusBadRequest
		case *ValidationError:
			httpStatus = http.StatusBadRequest
		case *service.DataStoreError:
			httpStatus = http.StatusInternalServerError
		default:
			httpStatus = http.StatusInternalServerError
		}
	}
	w.WriteHeader(httpStatus)
	if encodeErr := json.NewEncoder(w).Encode(err.Error()); encodeErr != nil {
		log.Error(encodeErr)
	}
}

// ReturnBodyJSON will return the body object as a JSON message
func ReturnBodyJSON(w http.ResponseWriter, body interface{}, httpStatus int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpStatus)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Error(err)
	}
}

// ReturnBlankJSON will return an empty JSON response
func ReturnBlankJSON(w http.ResponseWriter, httpStatus int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpStatus)
}
