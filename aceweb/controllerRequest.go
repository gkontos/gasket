package aceweb

import (
	"encoding/json"

	"io"
	"io/ioutil"
	"net/http"

	log "github.com/gkontos/gasket/acelog"
)

// ParseJsonRequest will parse the request body and return objects in the type of val interface{}
func ParseJsonRequest(r *http.Request, val interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		log.Error(err)
		return &RequestParseError{Message: "Unable to read message body", Err: err}
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, &val)
	if err != nil {
		log.Error(err)
		return &RequestParseError{Message: "Unable to parse request", Err: err}
	}
	return nil
}

// GetRequestBody will get the request body as a slice of byte
func GetRequestBody(r *http.Request) ([]byte, error) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, &RequestParseError{Message: "Unable to read message body", Err: err}
	}
	if err = r.Body.Close(); err != nil {
		return nil, &RequestParseError{Message: "Unable to close message body", Err: err}
	}

	return body, nil
}
