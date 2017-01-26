package model

import (
	"encoding/json"
	"fmt"

	"net/http"

	log "github.com/gkontos/gasket/acelog"
	"github.com/cayleygraph/cayley/quad"
)

// Node is the go struct for mapping an object to quads
type Node struct {
	ID         quad.IRI   `json:"id,omitempty"`
	Name       string     `json:"name"`
	Label      quad.Value `json:"label,omitempty"`
	Properties map[string]quad.Value
}

// NodeProperty is a utility struct for mapping quads
type NodeProperty struct {
	Predicate quad.IRI   `json:"predicate"`
	Object    quad.Value `json:"object"`
	Label     quad.Value `json:"label,omitempty"`
}

// PropertyValue is needed for testing, cannot initialize the map[string]quad.Value directly.  Might be a good replacement for the map?
type PropertyValue struct {
	Key   string
	Value quad.Value
}

// UnmarshalJSON will read a JSON object into a node.  name, id, label will be extracted other variables will be added to a map
func (np *Node) UnmarshalJSON(data []byte) error {

	var aux map[string]interface{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if _, ok := aux["id"]; ok {
		np.ID = quad.IRI(aux["id"].(string))
		delete(aux, "id")
	}
	if _, ok := aux["label"]; ok && aux["label"] != nil { // check if the key exists and is not nil
		np.Label = quad.String(aux["label"].(string))
		delete(aux, "label")
	} else if _, ok := aux["label"]; ok { // if the key exists, but it is nil; remove the key
		delete(aux, "label")
	} else {
		np.Label = (quad.Value)(nil)
	}

	if _, ok := aux["name"]; ok && aux["name"] != "" {
		log.Debug("setting name")

		np.Name = aux["name"].(string)
		delete(aux, "name")
	} else if _, ok := aux["name"]; ok {
		delete(aux, "name")
	} else {
		np.Name = ""
	}
	np.Properties = mapToValue(aux)

	return nil

}

// MarshalJSON will create a JSON object from a node.  The properties map will be expanded
func (np Node) MarshalJSON() ([]byte, error) {

	var err error

	var labelJSON []byte
	if _, ok := np.Label.(quad.Value); ok {
		if labelJSON, err = json.Marshal(np.Label.Native()); err != nil {
			return nil, err
		}
	} else {
		labelJSON = []byte(`null`)
	}

	var propertiesJSON string
	if propertiesJSON, err = getPropertiesJSONString(np.Properties); err != nil {
		return nil, err
	}

	jsonString := fmt.Sprintf("{\"id\":\"%s\",\"label\":%s,\"name\":\"%s\"%s}", UnEscapeIRI(np.ID), labelJSON, np.Name, propertiesJSON)

	return []byte(jsonString), nil
}

// UnmarshalJSON create a string from a NodeProperty struct
// http://choly.ca/post/go-json-marshalling/
func (np *NodeProperty) UnmarshalJSON(data []byte) error {

	type Alias NodeProperty

	aux := &struct {
		Predicate string `json:"predicate"`
		Object    string `json:"object"`
		Label     string `json:"label"`
		*Alias
	}{

		Alias: (*Alias)(np),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	np.Predicate = quad.IRI(aux.Predicate)
	np.Object = quad.Raw(aux.Object)
	np.Label = quad.Raw(aux.Label)

	return nil

}

func NewProperties(properties ...PropertyValue) map[string]quad.Value {
	propertyMap := make(map[string]quad.Value)
	for _, prop := range properties {
		propertyMap[prop.Key] = prop.Value
	}
	return propertyMap
}

// InputValidation should be a common interface used within the controllers
// https://husobee.github.io/golang/validation/2016/01/08/input-validation.html
type InputValidation interface {
	Validate(r *http.Request) error
}
