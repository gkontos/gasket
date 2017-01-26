package aceweb

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	log "github.com/gkontos/gasket/acelog"
	internal "github.com/gkontos/gasket/aceweb/internal"
	"github.com/gkontos/gasket/model"
	"github.com/cayleygraph/cayley/quad"
	"github.com/stretchr/testify/assert"
)

func TestMetadataCreateController(t *testing.T) {
	layout := "2006-01-02"
	timeValue, err := time.Parse(layout, "1964-02-23")
	if err != nil {
		log.Debug("error parsing time", err)
	}
	log.Debug(timeValue)
	tests := []internal.ControllerTestCase{
		{
			Description:    "Junk Input",
			Url:            "/metadata",
			Body:           []byte(`{"title":"Buy cheese and bread for breakfast."}`),
			ExpectedObject: &model.Metadata{}, // a blank expected object will not run tests for the property values
			ExpectedCode:   http.StatusCreated,
		}, {
			Description: "OK",
			Url:         "/metadata",
			Body: []byte(`{"relationId" : "klmnopqrst001",
								  "source" : "coffee shop",
								  "agree" : false,
								  "since" : "1964-02-23",
								  "popularity" : 15}`),
			ExpectedObject: &model.Metadata{
				RelationID: "klmnopqrst001",
				Properties: model.NewProperties(
					model.PropertyValue{Key: "source", Value: quad.Raw("coffee shop")},
					model.PropertyValue{Key: "agree", Value: quad.Bool(false)},
					model.PropertyValue{Key: "since", Value: quad.Raw("1964-02-23")},
					model.PropertyValue{Key: "popularity", Value: quad.Float(15)},
				)},
			ExpectedCode: http.StatusCreated,
		},
	}

	internal.RunControllerTests(t, tests, "POST", http.HandlerFunc(MetadataAdd),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {
			log.Debug("running test case : ", tc.Description)
			log.Debug("received : ", string(body))
			assert := assert.New(t)
			var metadata model.Metadata
			err := json.Unmarshal(body, &metadata)
			assert.NoError(err)
			expectedMetadata := tc.ExpectedObject.(*model.Metadata)
			assert.Equal(expectedMetadata.RelationID, metadata.RelationID, tc.Description+" -id")
			for key, expectedValue := range expectedMetadata.Properties {
				if assert.NotNil(metadata.Properties[key], key+" is nil") {
					assert.Equal(expectedValue, metadata.Properties[key], tc.Description+" key="+key)
				}
			}
		})
}

func TestMetadataDelete(t *testing.T) {
	idExists := "zyx987654321"
	//	idDoesNotExist := "IWillNotBeFound"
	tests := []internal.ControllerTestCase{
		{
			Description:    "Metadata exists",
			Url:            "/metadata/" + idExists,
			Body:           []byte(``),
			ExpectedObject: &model.Metadata{}, // a blank expected object will not run tests for the property values
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

	internal.RunControllerTests(t, tests, "DELETE", http.HandlerFunc(MetadataDelete),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {

		})
}

func TestMetadataGetController(t *testing.T) {
	idExists := "zyx987654321"
	idDoesNotExist := "IWillNotBeFound"
	tests := []internal.ControllerTestCase{
		{
			Description: "Node exists",
			RouteUrl:    "/metadata/{metadataid}",
			Url:         "/metadata/" + idExists,
			Body:        []byte(``),
			ExpectedObject: &model.Metadata{
				RelationID: quad.IRI("abcdefghij001"),
				Properties: model.NewProperties(
					model.PropertyValue{Key: "source", Value: quad.Raw("interweb.com")},
					model.PropertyValue{Key: "popularity", Value: quad.Raw("95010")},
					model.PropertyValue{Key: "dateCreated", Value: quad.Raw("1989-10-01")},
					model.PropertyValue{Key: "agree", Value: quad.Raw("false")},
				)},
			ExpectedCode: http.StatusOK,
		}, {
			Description:    "Does not exist",
			RouteUrl:       "/metadata/{metadataid}",
			Url:            "/metadata/" + idDoesNotExist,
			Body:           nil,
			ExpectedObject: nil,
			ExpectedCode:   http.StatusNotFound,
		},
	}

	internal.RunControllerTests(t, tests, "GET", http.HandlerFunc(MetadataGet),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {
			log.Debug("running test case : ", tc.Description)
			log.Debug("received : " + string(body))
			assert := assert.New(t)
			if tc.Body != nil && tc.ExpectedObject != nil {
				var metadata model.Metadata
				err := json.Unmarshal(body, &metadata)
				assert.NoError(err)
				expectedMetadata := tc.ExpectedObject.(*model.Metadata)
				assert.Equal(expectedMetadata.RelationID, metadata.RelationID, tc.Description+" -label")
				for key, expectedValue := range expectedMetadata.Properties {
					assert.Equal(expectedValue, metadata.Properties[key], tc.Description+" key="+key)
				}
			}
		})
}

func TestMetadataPutController(t *testing.T) {
	idExists := "zyx987654321"
	//	idDoesNotExist := "IWillNotBeFound"
	tests := []internal.ControllerTestCase{
		{
			Description: "Metadata Update and Add",
			RouteUrl:    "/metadata/{metadataid}",
			Url:         "/metadata/" + idExists,
			Body: []byte(`{"relationId" : "abcdefghij001",
								  "source" : "coffee shop",
								  "agree" : false,
								  "since" : "1964-02-23",
								  "popularity" : 15}`),
			ExpectedObject: &model.Metadata{
				RelationID: quad.IRI("abcdefghij001"),
				Properties: model.NewProperties(
					model.PropertyValue{Key: "source", Value: quad.Raw("coffee shop")},
					model.PropertyValue{Key: "popularity", Value: quad.Float(15)},
					model.PropertyValue{Key: "since", Value: quad.Raw("1964-02-23")},
					model.PropertyValue{Key: "dateCreated", Value: quad.Raw("1989-10-01")},
					model.PropertyValue{Key: "agree", Value: quad.Bool(false)},
				)},
			ExpectedCode: http.StatusOK,
		}, {
			Description: "Metadata Add",
			RouteUrl:    "/metadata/{metadataid}",
			Url:         "/metadata/yx9876543210",
			Body:        []byte(`{"has facial hair" : "fu man chu"}`),
			ExpectedObject: &model.Metadata{
				RelationID: quad.IRI("abcdefghij001"),
				Properties: model.NewProperties(
					model.PropertyValue{Key: "source", Value: quad.Raw("grandma")},
					model.PropertyValue{Key: "popularity", Value: quad.Raw("1")},
					model.PropertyValue{Key: "dateCreated", Value: quad.Raw("2001-05-01")},
					model.PropertyValue{Key: "has facial hair", Value: quad.Raw("fu man chu")},
					model.PropertyValue{Key: "agree", Value: quad.Raw("true")},
				)},
			ExpectedCode: http.StatusOK,
		},
		//		, {
		//			Description: "Does not exist",
		//			RouteUrl:    "/metadata/{metadataid}",
		//			Url:         "/metadata/" + idDoesNotExist,
		//			Body: []byte(`{"label" : "test",
		//								  "name" : "node create test",
		//								  "type" : "acedfs:process",
		//								  "color" : "orange",
		//								  "amount" : 11.11}`),
		//			ExpectedObject: nil,
		//			ExpectedCode:   http.StatusNotFound,
		//		},
	}

	internal.RunControllerTests(t, tests, "PUT", http.HandlerFunc(MetadataUpdate),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {
			log.Debug("running test case : ", tc.Description)
			log.Debug("received : " + string(body))
			assert := assert.New(t)
			if tc.Body != nil && tc.ExpectedObject != nil {
				var metadata model.Metadata
				err := json.Unmarshal(body, &metadata)
				assert.NoError(err)
				expectedMetadata := tc.ExpectedObject.(*model.Metadata)
				assert.Equal(expectedMetadata.RelationID, metadata.RelationID, tc.Description+" -label")
				for key, expectedValue := range expectedMetadata.Properties {
					assert.Equal(expectedValue, metadata.Properties[key], tc.Description+" key="+key)
				}
			}
		})
}
