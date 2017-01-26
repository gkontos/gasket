package aceservice

import (
	"fmt"

	"github.com/gkontos/gasket/model"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/quad"
	"github.com/pborman/uuid"
)

// DeleteByID removes all nodes with the value of 'subject'.  This value may be the label, object, subject or predicate
func DeleteByID(subject string) error {

	return store.RemoveNode(store.ValueOf(quad.IRI(subject)))
}

// NodeToNodeProperties will return the properties map of the node as a list of NodeProperty
func NodeToNodeProperties(node model.Node) ([]model.NodeProperty, error) {
	var nodeProperties []model.NodeProperty
	var err error
	for key, value := range node.Properties {
		objectVal, ok := quad.AsValue(value)
		if ok {
			nodeProperty := model.NodeProperty{
				Predicate: quad.IRI(key),
				Object:    objectVal,
				Label:     node.Label,
			}
			nodeProperties = append(nodeProperties, nodeProperty)
		} else {
			err = fmt.Errorf("Unable to parse property : %s", value)
		}
	}
	nameProperty := model.NodeProperty{
		Predicate: model.NamePredicate,
		Object:    quad.String(node.Name),
		Label:     node.Label,
	}
	nodeProperties = append(nodeProperties, nameProperty)
	return nodeProperties, err
}

// AddNode will save a node and the node properties as quads to the data store
func AddNode(node model.Node) ([]quad.Quad, error) {
	var quadList []quad.Quad
	nodeProperties, parseErr := NodeToNodeProperties(node)
	nodeID := uuid.NewUUID()

	if parseErr != nil {
		return nil, parseErr
	}
	var err error
	tx := cayley.NewTransaction()
	for _, nodeProperty := range nodeProperties {
		propertyQuad := quad.Make(quad.IRI(nodeID.String()),
			nodeProperty.Predicate,
			nodeProperty.Object,
			nodeProperty.Label)
		tx.AddQuad(propertyQuad)
		quadList = append(quadList, propertyQuad)
	}
	err = store.ApplyTransaction(tx)
	if err != nil {
		return nil, &DataStoreError{Message: "Error saving data", Err: err}
	}
	return quadList, err
}

// UpdateNode will add or update any properties of the node
func UpdateNode(node model.Node) ([]quad.Quad, error) {

	var quadList []quad.Quad
	var saveErr error
	nodeProperties, parseErr := NodeToNodeProperties(node)
	if parseErr != nil {
		return nil, parseErr
	}
	tx := cayley.NewTransaction()
	for _, nodeProperty := range nodeProperties {

		if nodeProperty.Predicate != model.NamePredicate && nodeProperty.Object.String() != "\"\"" { // don't process a missing name attribute
			propertyQuad := quad.Make(node.ID,
				nodeProperty.Predicate,
				nodeProperty.Object,
				nodeProperty.Label)
			AddOrUpdateAsTransaction(tx, propertyQuad)
			quadList = append(quadList, propertyQuad)
		}
	}

	err := store.ApplyTransaction(tx)
	if err != nil {
		saveErr = &DataStoreError{Message: "Error updating data", Err: err}
		return nil, saveErr
	}
	return quadList, nil
}
