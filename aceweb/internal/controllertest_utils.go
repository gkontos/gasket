package aceweb

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gkontos/gasket/aceservice"
	"github.com/cayleygraph/cayley"
	_ "github.com/cayleygraph/cayley/graph/memstore"
	"github.com/cayleygraph/cayley/quad"
	"github.com/cayleygraph/cayley/quad/cquads"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func MakeTestStore(t testing.TB) *cayley.Handle {
	// TODO we don't want to load this each time a test is run
	simpleGraph := LoadGraph(t, "./../data/testdata-3.nq")
	h, err := cayley.NewGraph("memstore", "", nil)

	if err != nil {
		t.Fatalf("Failed to setup test datastore")
	}

	for _, q := range simpleGraph {
		h.AddQuad(q)
	}
	aceservice.SetStore(h)
	return h
}

type ControllerTestCase struct {
	Description    string
	Url            string
	RouteUrl       string // must be set for gorilla/mux to pick up url parameters (otherwise not needed)
	Body           []byte
	ExpectedObject interface{}
	ExpectedCode   int
}

type ResponseTest func(t *testing.T, body []byte, testCase ControllerTestCase)

func RunControllerTests(t *testing.T, tests []ControllerTestCase, httpMethod string, handler http.HandlerFunc, responseTest ResponseTest) {
	MakeTestStore(t)
	assert := assert.New(t)

	for _, tc := range tests {
		var u bytes.Buffer
		//		u.WriteString(string(ts.URL))
		u.WriteString(tc.Url)

		fmt.Println(u.String())
		req, err := http.NewRequest(httpMethod, u.String(), bytes.NewBuffer(tc.Body))
		req.Header.Set("Content-Type", "application/json")

		assert.NoError(err)
		resp := httptest.NewRecorder()
		if tc.RouteUrl == "" {
			tc.RouteUrl = tc.Url
		}
		getRouter(tc.RouteUrl, handler).ServeHTTP(resp, req)

		b, err := ioutil.ReadAll(resp.Body)
		assert.NoError(err)
		assert.Equal(tc.ExpectedCode, resp.Code, tc.Description)
		responseTest(t, b, tc)

	}
}

func getRouter(routeUrl string, handler http.HandlerFunc) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc(routeUrl, handler)

	return r
}

func LoadGraph(t testing.TB, path string) []quad.Quad {
	var r io.Reader
	var simpleGraph []quad.Quad
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open %q: %v", path, err)
	}
	defer f.Close()
	r = f

	dec := cquads.NewDecoder(r)
	for q1, err := dec.Unmarshal(); err == nil; q1, err = dec.Unmarshal() {
		simpleGraph = append(simpleGraph, q1)
	}
	if err != nil {
		t.Fatalf("Failed to Unmarshal: %v", err)
	}
	return simpleGraph

}
