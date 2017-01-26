package aceweb

import (
	"net/http"
	"testing"

	"encoding/json"

	log "github.com/gkontos/gasket/acelog"
	internal "github.com/gkontos/gasket/aceweb/internal"
	"github.com/gkontos/gasket/model"
	"github.com/cayleygraph/cayley/quad"
	"github.com/stretchr/testify/assert"
)

func TestNodeCreateController(t *testing.T) {

	tests := []internal.ControllerTestCase{
		{
			Description:    "Junk Input",
			Url:            "/nodes",
			Body:           []byte(`{"title":"Buy cheese and bread for breakfast."}`),
			ExpectedObject: &model.Node{}, // a blank expected object will not run tests for the property values
			ExpectedCode:   http.StatusCreated,
		}, {
			Description: "OK",
			Url:         "/nodes",
			Body: []byte(`{"label" : "test",
								  "name" : "node create test",
								  "type" : "acedfs:process",
								  "color" : "orange",
								  "amount" : 11.11}`),
			ExpectedObject: &model.Node{
				Label: quad.String("test"),
				Name:  "node create test",
				Properties: model.NewProperties(
					model.PropertyValue{Key: "type", Value: quad.Raw("acedfs:process")},
					model.PropertyValue{Key: "color", Value: quad.Raw("orange")},
					model.PropertyValue{Key: "amount", Value: quad.Float(11.11)},
				)},
			ExpectedCode: http.StatusCreated,
		},
	}

	internal.RunControllerTests(t, tests, "POST", http.HandlerFunc(NodeCreate),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {
			log.Debug("running test case : ", tc.Description)
			assert := assert.New(t)
			var node model.Node
			err := json.Unmarshal(body, &node)
			assert.NoError(err)
			expectedNode := tc.ExpectedObject.(*model.Node)

			assert.Equal(expectedNode.Label, node.Label, tc.Description+" -label")
			assert.Equal(expectedNode.Name, node.Name, tc.Description+" -name")
			for key, expectedValue := range expectedNode.Properties {
				assert.Equal(expectedValue, node.Properties[key], tc.Description+" key="+key)
			}
		})
}

func TestNodeDeleteController(t *testing.T) {
	idExists := "123456789"
	//	idDoesNotExist := "IWillNotBeFound"
	tests := []internal.ControllerTestCase{
		{
			Description:    "Node exists",
			Url:            "/nodes/" + idExists,
			Body:           []byte(``),
			ExpectedObject: &model.Node{}, // a blank expected object will not run tests for the property values
			ExpectedCode:   http.StatusNoContent,
		},
		//		{
		//			Description:    "Does not exist",
		//			Url:            "/nodes/" + idDoesNotExist,
		//			Body:           []byte(``),
		//			ExpectedObject: &model.Node{},
		//			ExpectedCode:   http.StatusNotFound,
		//		},
	}

	internal.RunControllerTests(t, tests, "DELETE", http.HandlerFunc(NodeDelete),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {

		})
}

func TestNodeGetController(t *testing.T) {
	idExists := "123456789"
	idDoesNotExist := "IWillNotBeFound"
	tests := []internal.ControllerTestCase{
		{
			Description: "Node exists",
			RouteUrl:    "/nodes/{id}",
			Url:         "/nodes/" + idExists,
			Body:        []byte(``),
			ExpectedObject: &model.Node{
				Label: quad.String("test"),
				Name:  "Shimmering Substance",
				Properties: model.NewProperties(
					model.PropertyValue{Key: "color", Value: quad.Raw("yellow")},
					model.PropertyValue{Key: "rdf:type", Value: quad.Raw("painting")},
					model.PropertyValue{Key: "style", Value: quad.Raw("abstract expressionist")},
				)},
			ExpectedCode: http.StatusOK,
		}, {
			Description:    "Does not exist",
			RouteUrl:       "/nodes/{id}",
			Url:            "/nodes/" + idDoesNotExist,
			Body:           nil,
			ExpectedObject: nil,
			ExpectedCode:   http.StatusNotFound,
		},
	}

	internal.RunControllerTests(t, tests, "GET", http.HandlerFunc(NodeGet),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {
			log.Debug("running test case : ", tc.Description)
			log.Debug("received : " + string(body))
			assert := assert.New(t)
			if tc.Body != nil && tc.ExpectedObject != nil {
				var node model.Node
				err := json.Unmarshal(body, &node)
				assert.NoError(err)
				expectedNode := tc.ExpectedObject.(*model.Node)
				assert.Equal(expectedNode.Label, node.Label, tc.Description+" -label")
				assert.Equal(expectedNode.Name, node.Name, tc.Description+" -name")
				for key, expectedValue := range expectedNode.Properties {
					assert.Equal(expectedValue, node.Properties[key], tc.Description+" key="+key)
				}
			}
		})
}

func TestNodePutController(t *testing.T) {
	idExists := "123456789"
	//	idDoesNotExist := "IWillNotBeFound"
	tests := []internal.ControllerTestCase{
		{
			Description: "Add Property",
			RouteUrl:    "/nodes/{id}",
			Url:         "/nodes/" + idExists,
			Body:        []byte(`{"year" : "1946"}`),
			ExpectedObject: &model.Node{
				Label: quad.String("test"),
				Name:  "Shimmering Substance",
				Properties: model.NewProperties(
					model.PropertyValue{Key: "color", Value: quad.Raw("yellow")},
					model.PropertyValue{Key: "rdf:type", Value: quad.Raw("painting")},
					model.PropertyValue{Key: "style", Value: quad.Raw("abstract expressionist")},
					model.PropertyValue{Key: "year", Value: quad.Raw("1946")},
				)},
			ExpectedCode: http.StatusOK,
		}, {
			Description: "Change Property",
			RouteUrl:    "/nodes/{id}",
			Url:         "/nodes/" + idExists,
			Body:        []byte(`{"color" : "yellow hues"}`),
			ExpectedObject: &model.Node{
				Label: quad.String("test"),
				Name:  "Shimmering Substance",
				Properties: model.NewProperties(
					model.PropertyValue{Key: "color", Value: quad.Raw("yellow hues")},
					model.PropertyValue{Key: "rdf:type", Value: quad.Raw("painting")},
					model.PropertyValue{Key: "style", Value: quad.Raw("abstract expressionist")},
				)},
			ExpectedCode: http.StatusOK,
		},
		//		{
		//			Description: "Does not exist",
		//			RouteUrl:    "/nodes/{id}",
		//			Url:         "/nodes/" + idDoesNotExist,
		//			Body: []byte(`{"label" : "test",
		//								  "name" : "node create test",
		//								  "type" : "acedfs:process",
		//								  "color" : "orange",
		//								  "amount" : 11.11}`),
		//			ExpectedObject: nil,
		//			ExpectedCode:   http.StatusNotFound,
		//		},
	}

	internal.RunControllerTests(t, tests, "PUT", http.HandlerFunc(NodeUpdate),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {
			log.Debug("running test case : ", tc.Description)
			log.Debug("received : " + string(body))
			assert := assert.New(t)
			if tc.Body != nil && tc.ExpectedObject != nil {
				var node model.Node
				err := json.Unmarshal(body, &node)
				assert.NoError(err)
				expectedNode := tc.ExpectedObject.(*model.Node)
				assert.Equal(expectedNode.Label, node.Label, tc.Description+" -label")
				assert.Equal(expectedNode.Name, node.Name, tc.Description+" -name")
				for key, expectedValue := range expectedNode.Properties {
					assert.Equal(expectedValue, node.Properties[key], tc.Description+" key="+key)
				}
			}
		})
}
