package model

import (
	"encoding/json"
	"fmt"

	"github.com/cayleygraph/cayley/quad"
)

type Relation struct {
	ID       quad.IRI   `json:"id,omitempty"`
	SourceID quad.IRI   `json:"sourceId"`
	Type     quad.IRI   `json:"type"`
	TargetID quad.IRI   `json:"targetId"`
	Label    quad.Value `json:"label,omitempty"`
}

//RelationQuad is a typed quad for interacting with cayley
type RelationQuad struct {
	Subject   quad.IRI    `json:"subject"`
	Predicate quad.IRI    `json:"predicate"`
	Object    quad.IRI    `json:"object"`
	Label     quad.String `json:"label,omitempty"`
}

func (m Relation) MarshalJSON() ([]byte, error) {

	var labelJSON []byte
	var err error
	if _, ok := m.Label.(quad.Value); ok {
		if labelJSON, err = json.Marshal(m.Label.Native()); err != nil {
			return nil, err
		}
	} else {
		labelJSON = []byte(`null`)
	}

	jsonString := fmt.Sprintf("{\"id\":\"%s\",\"sourceId\":\"%s\",\"type\":\"%s\",\"targetId\":\"%s\",\"label\":%s}",
		UnEscapeIRI(m.ID),
		UnEscapeIRI(m.SourceID),
		UnEscapeIRI(m.Type),
		UnEscapeIRI(m.TargetID),
		labelJSON)

	return []byte(jsonString), nil
}
