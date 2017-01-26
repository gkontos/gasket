package aceservice

import (
	"fmt"

	"reflect"

	"github.com/gkontos/gasket/model"
	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/quad"
)

// QuadListToNode will map a list of quads to a Node.  If the quads do not share a node id as subject, an error will be thrown
func QuadListToNode(quads []quad.Quad) (model.Node, error) {
	var node model.Node

	node.Properties = make(map[string]quad.Value)

	var err error
	for _, q := range quads {
		if node.ID != "" {
			if q.Subject != node.ID {
				err = fmt.Errorf("Quads are not all from the same node.")
				break
			}
		} else {
			node.ID = q.Subject.(quad.IRI)
		}
		if _, ok := node.Label.(quad.Value); !ok {
			// TODO ?multiple labels for nodes?
			node.Label = q.Label
		}

		if q.Predicate == model.NamePredicate {
			node.Name = q.Object.Native().(string)
		} else {
			if reflect.TypeOf(q.Predicate).String() == "quad.IRI" {
				node.Properties[model.UnEscapeIRI(q.Predicate.(quad.IRI))] = q.Object
			} else {
				node.Properties[q.Predicate.String()] = q.Object
			}
		}
	}
	return node, err
}

// QuadListToNode will map a list of quads to a Node.  If the quads do not share a node id as subject, an error will be thrown
func QuadListToMetadata(quads []quad.Quad) (model.Metadata, error) {
	var metadata model.Metadata

	metadata.Properties = make(map[string]quad.Value)

	var err error
	for _, q := range quads {
		if q.Predicate == model.MetaidPredicate {

			if metadata.ID.String() == "<>" {
				metadata.ID = q.Object.(quad.IRI)
			}
			metadata.RelationID = q.Subject.(quad.IRI)

		} else {
			if metadata.ID != "" && q.Subject != metadata.ID {
				err = fmt.Errorf("Quads are not all from the same node for metadataid=%s", metadata.ID)
				break
			}

			if reflect.TypeOf(q.Predicate).String() == "quad.IRI" {

				metadata.Properties[model.UnEscapeIRI(q.Predicate.(quad.IRI))] = q.Object
			} else {
				metadata.Properties[q.Predicate.String()] = q.Object
			}
		}
	}
	return metadata, err
}

// AddQuad will save the given quad to the store
func AddQuad(q quad.Quad) error {

	err := store.AddQuad(q)
	if err != nil {
		dberr := &DataStoreError{Message: "Unable to save to datastore", Err: err}
		return dberr
	}
	return nil
}

// DeleteQuad will delete the given quad from the store
func DeleteQuad(q quad.Quad) error {

	err := store.RemoveQuad(q)
	if err != nil {
		dberr := &DataStoreError{Message: "Unable to delete from datastore", Err: err}
		return dberr
	}
	return nil
}

// AddOrUpdate quad will add the new quad to the store
// if a quad is found matching the subject and prediate, the existing quad will be deleted
func AddOrUpdate(q quad.Quad) error {

	tx := cayley.NewTransaction()
	AddOrUpdateAsTransaction(tx, q)
	err := store.ApplyTransaction(tx)
	if err != nil {
		dberr := &DataStoreError{Message: "Transaction error", Err: err}
		return dberr
	}
	return nil
}

func AddOrUpdateAsTransaction(tx *graph.Transaction, q quad.Quad) {

	for _, direction := range []quad.Direction{quad.Subject} {
		it := store.QuadIterator(direction, store.ValueOf(q.Subject))

		for it.Next() {

			foundQuad := store.Quad(it.Result())

			if foundQuad.Predicate == q.Predicate {
				tx.RemoveQuad(foundQuad)
			}
		}
	}
	tx.AddQuad(q)

}

// GetQuads will return all quads will subject or objects containing the parameter {subject}
func GetQuads(subject string) []quad.Quad {

	var quadList []quad.Quad

	// see writer/single.go for an example function
	for _, direction := range []quad.Direction{quad.Subject, quad.Object} {
		it := store.QuadIterator(direction, store.ValueOf(quad.IRI(subject)))
		for it.Next() {

			quadList = append(quadList, store.Quad(it.Result()))

		}
		it.Close()
	}
	return quadList
}

// GetQuadsBySubject will return all quads will subject containing the parameter {subject}
func GetQuadsBySubject(subject string) []quad.Quad {

	var quadList []quad.Quad

	// see writer/single.go for an example function
	for _, direction := range []quad.Direction{quad.Subject} {
		it := store.QuadIterator(direction, store.ValueOf(quad.IRI(subject)))
		for it.Next() {

			quadList = append(quadList, store.Quad(it.Result()))

		}
		it.Close()
	}
	return quadList
}
