package aceweb

import (
	"encoding/json"
	"net/http"
	"testing"

	log "github.com/gkontos/gasket/acelog"
	internal "github.com/gkontos/gasket/aceweb/internal"
	"github.com/gkontos/gasket/model"
	_ "github.com/cayleygraph/cayley/graph/memstore"
	"github.com/cayleygraph/cayley/quad"
	"github.com/stretchr/testify/assert"
)

func TestRelationCreateController(t *testing.T) {
	tests := []internal.ControllerTestCase{
		{
			Description: "Basic Input",
			Url:         "/relations",
			Body:        []byte(`{"sourceid":"123456789", "type":"pavedthewayfor","targetid":"234567890"}`),
			ExpectedObject: &model.Relation{
				SourceID: "123456789",
				Type:     "pavedthewayfor",
				TargetID: "234567890",
			},
			ExpectedCode: http.StatusCreated,
		},
		{
			Description:    "Id does not exist",
			Url:            "/relations",
			Body:           []byte(`{"sourceid":"noid", "type":"pavedthewayfor","targetid":"234567890"}`),
			ExpectedObject: nil,
			ExpectedCode:   http.StatusCreated,
		},
	}

	internal.RunControllerTests(t, tests, "POST", http.HandlerFunc(RelationCreate),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {
			log.Debug("running test case : ", tc.Description)
			assert := assert.New(t)
			var relation model.Relation
			err := json.Unmarshal(body, &relation)
			assert.NoError(err)
			if tc.ExpectedObject != nil {
				expected := tc.ExpectedObject.(*model.Relation)
				assert.Equal(expected.SourceID, relation.SourceID, tc.Description)
				assert.Equal(expected.TargetID, relation.TargetID, tc.Description)
				assert.Equal(expected.Type, relation.Type, tc.Description)
			}

		})
}

func TestRelationDelete(t *testing.T) {
	idExists := "abcdefghij001"
	//	idDoesNotExist := "IWillNotBeFound"
	tests := []internal.ControllerTestCase{
		{
			Description:    "Node exists",
			Url:            "/relations/" + idExists,
			Body:           []byte(``),
			ExpectedObject: &model.Relation{}, // a blank expected object will not run tests for the property values
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

	internal.RunControllerTests(t, tests, "DELETE", http.HandlerFunc(RelationDelete),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {

		})
}

func TestQuadMarshal(t *testing.T) {
	assert := assert.New(t)
	q := quad.Quad{
		Subject:   quad.IRI("testsubject"),
		Predicate: quad.IRI("testpred"),
		Object:    quad.IRI("testobj"),
		Label:     quad.String("testlabel"),
	}

	quadJSON, err := json.Marshal(q)
	assert.NoError(err)

	var relationQuad quad.Quad
	err = json.Unmarshal(quadJSON, &relationQuad)
	assert.NoError(err)

	quadJSON, err = json.Marshal(q)

}

func TestRelationGetController(t *testing.T) {
	idExists := "abcdefghij001"
	idDoesNotExist := "IWillNotBeFound"
	tests := []internal.ControllerTestCase{
		{
			Description: "Node exists",
			RouteUrl:    "/relations/{id}",
			Url:         "/relations/" + idExists,
			Body:        []byte(``),
			ExpectedObject: &model.Relation{
				ID:       quad.IRI("abcdefghij001"),
				SourceID: quad.IRI("123456789"),
				Type:     quad.IRI("similarto"),
				TargetID: quad.IRI("234567890"),
				Label:    nil,
			},
			ExpectedCode: http.StatusOK,
		}, {
			Description:    "Does not exist",
			RouteUrl:       "/relations/{id}",
			Url:            "/relations/" + idDoesNotExist,
			Body:           nil,
			ExpectedObject: nil,
			ExpectedCode:   http.StatusNotFound,
		},
	}

	internal.RunControllerTests(t, tests, "GET", http.HandlerFunc(RelationGet),
		func(t *testing.T, body []byte, tc internal.ControllerTestCase) {
			log.Debug("running test case : ", tc.Description)
			assert := assert.New(t)
			if tc.Body != nil && tc.ExpectedObject != nil {
				var relation model.Relation
				err := json.Unmarshal(body, &relation)
				assert.NoError(err)
				expectedRelation := tc.ExpectedObject.(*model.Relation)

				assert.Equal(expectedRelation.ID, relation.ID, tc.Description+" -id")
				assert.Equal(expectedRelation.SourceID, relation.SourceID, tc.Description+" -source")
				assert.Equal(expectedRelation.Type, relation.Type, tc.Description+" -type")
				assert.Equal(expectedRelation.TargetID, relation.TargetID, tc.Description+" -target")
				assert.Equal(expectedRelation.Label, relation.Label, tc.Description+" -label")
			}
		})
}
