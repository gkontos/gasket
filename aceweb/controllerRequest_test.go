package aceweb

import (
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/gkontos/gasket/model"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Id uint32 `json:"id"`

	Username string `json:"username"`

	MoneyBalance uint32 `json:"balance"`
}

func TestSingleObject(t *testing.T) {

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	userJson := `{"username": "dennis", "balance": 200}`
	reader := strings.NewReader(userJson)
	req, err := http.NewRequest("GET", "/health-check", reader)
	if err != nil {
		t.Fatal(err)
	}
	var object User
	ParseJsonRequest(req, &object)
	if object.Username != "dennis" {
		t.Errorf("username is %s expected dennis", object.Username)
	}
}

func TestArrayObject(t *testing.T) {

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	userJson := `[{"username": "dennis", "balance": 200},{"username": "bob", "balance": 5000}]`
	reader := strings.NewReader(userJson)
	req, err := http.NewRequest("GET", "/health-check", reader)
	if err != nil {
		t.Fatal(err)
	}
	var object []User
	ParseJsonRequest(req, &object)
	if len(object) != 2 {
		t.Errorf("not all objects parsed")
	}

	assert.True(t, reflect.TypeOf(object).String() == "[]aceweb.User")
	for _, user := range object {
		assert.True(t, user.Username != "")
	}

}

func TestBadArrayObject(t *testing.T) {
	assert := assert.New(t)
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	userJson := `{"username": "dennis", "balance": 200},{"username": "bob", "balance": 5000}]`
	reader := strings.NewReader(userJson)
	req, err := http.NewRequest("GET", "/health-check", reader)
	if err != nil {
		t.Fatal(err)
	}
	var object []User
	err = ParseJsonRequest(req, &object)
	assert.Error(err)

}

func TestNodeParse(t *testing.T) {
	assert := assert.New(t)
	userJson := `{"label" : "test",
								  "name" : "node create test",
								  "type" : "acedfs:process",
								  "color" : "orange",
								  "amount" : 11.11}`
	reader := strings.NewReader(userJson)
	req, err := http.NewRequest("GET", "/health-check", reader)
	if err != nil {
		t.Fatal(err)
	}
	var object model.Node
	err = ParseJsonRequest(req, &object)
	assert.NoError(err)
	assert.Equal("test", object.Label.Native(), "parsed label property")
}

func TestNodeParseRandom(t *testing.T) {
	assert := assert.New(t)
	userJson := `{"title":"Buy cheese and bread for breakfast."}`
	reader := strings.NewReader(userJson)
	req, err := http.NewRequest("GET", "/health-check", reader)
	if err != nil {
		t.Fatal(err)
	}
	var object model.Node
	err = ParseJsonRequest(req, &object)
	assert.NoError(err)
	assert.Equal("Buy cheese and bread for breakfast.", object.Properties["title"].String(), "parsed property")
}
