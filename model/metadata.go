package model

import (
	"encoding/json"
	"fmt"

	"github.com/cayleygraph/cayley/quad"
)

// Metadata is the go type for relation metadata
type Metadata struct {
	ID         quad.IRI `json:"id,omitempty"`
	RelationID quad.IRI `json:"relationId"`
	Properties map[string]quad.Value
}

// UnmarshalJSON will create a JSON object from metadata object.  The properties map will be expanded
func (m *Metadata) UnmarshalJSON(data []byte) error {

	var aux map[string]interface{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if _, ok := aux["id"]; ok {
		m.ID = quad.IRI(aux["id"].(string))
		delete(aux, "id")
	}

	if _, ok := aux["relationId"]; ok && aux["relationId"] != "" {
		m.RelationID = quad.IRI(aux["relationId"].(string))
		delete(aux, "relationId")
	} else if _, ok := aux["relationId"]; ok {
		delete(aux, "relationId")
	} else {
		m.RelationID = quad.IRI("")
	}
	m.Properties = mapToValue(aux)

	return nil

}

func (m Metadata) MarshalJSON() ([]byte, error) {

	var err error

	var propertiesJSON string

	if propertiesJSON, err = getPropertiesJSONString(m.Properties); err != nil {
		return nil, err
	}

	jsonString := fmt.Sprintf("{\"id\":\"%s\",\"relationId\":\"%s\"%s}", UnEscapeIRI(m.ID), UnEscapeIRI(m.RelationID), propertiesJSON)

	return []byte(jsonString), nil
}
