package model

import (
	"bytes"
	"encoding/json"

	"fmt"

	"github.com/cayleygraph/cayley/quad"
)

const RelationidPredicate = quad.IRI("hasRelationId")
const MetaidPredicate = quad.IRI("hasMetaId")
const NamePredicate = quad.IRI("schema:name")

//mapToValue will return a map of cayley typed values based on a json map input
func mapToValue(jsonmap map[string]interface{}) map[string]quad.Value {
	propmap := make(map[string]quad.Value)
	for key, value := range jsonmap {
		switch value.(type) {
		case string:
			propmap[key] = quad.Raw(value.(string))
		case float64:
			propmap[key], _ = quad.AsValue(value.(float64))
		case int:
			propmap[key], _ = quad.AsValue(value.(int))
		default:
			propmap[key], _ = quad.AsValue(value)
		}
		// TODO check for time
	}
	return propmap
}

func getPropertiesJSONString(props map[string]quad.Value) (string, error) {
	var propertiesJSON bytes.Buffer

	var valueJSON []byte
	var err error
	for key, value := range props {

		if valueJSON, err = json.Marshal(value.Native()); err != nil {
			return "", err
		}
		propertiesJSON.WriteString("," + fmt.Sprintf("\"%s\":%s", key, valueJSON))
	}
	return propertiesJSON.String(), nil
}

// UnEscapeIRI removes the `<>` from an IRI string representation.  would be nice as a method on IRI, would require wrapping the IRI struct
func UnEscapeIRI(iri quad.Value) string {
	s := iri.String()
	s = s[1 : len(s)-1]
	return s

}
