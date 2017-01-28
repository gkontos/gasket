package aceservice

import (
	"encoding/json"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/iterator"
	"github.com/cayleygraph/cayley/quad"
	log "github.com/gkontos/gasket/acelog"
	"github.com/gkontos/gasket/model"
	"github.com/pborman/uuid"
)

// GetRelationshipSubject will return the subject of a relation as a quad.
func GetRelationshipSubject(baseQuad quad.Quad) (string, error) {
	bytes, err := json.Marshal(baseQuad)

	return string(bytes), err
}

// GetRelation will return the relation for an ID
func GetRelation(ID string) model.Relation {

	var relation model.Relation
	var relationQuad quad.Quad
	var foundQuad quad.Quad

	it, _ := iterator.NewAnd(
		store,
		store.QuadIterator(quad.Object, store.ValueOf(quad.IRI(ID))),
		store.QuadIterator(quad.Predicate, store.ValueOf(model.RelationidPredicate)),
	).Optimize()
	defer it.Close()

	it, _ = store.OptimizeIterator(it)

	for it.Next() {
		// we are only expecting a single quad with a specific id, so once found, break
		foundQuad = store.Quad(it.Result())

		break
	}
	if _, ok := foundQuad.Subject.(quad.Value); ok {

		if err := json.Unmarshal([]byte(foundQuad.Subject.Native().(string)), &relationQuad); err == nil {

			relation.ID = quad.IRI(ID)
			relation.SourceID = quad.IRI(model.UnEscapeIRI(relationQuad.Subject))
			relation.Type = quad.IRI(model.UnEscapeIRI(relationQuad.Predicate))
			relation.TargetID = quad.IRI(model.UnEscapeIRI(relationQuad.Object))
			relation.Label = relationQuad.Label
		} else {
			log.Error(err)
		}
	}
	return relation
}

// DeleteByRelationID will Delete the the relation quad, and any metadata quads for the given relationid
func DeleteByRelationID(ID string) error {
	var deleteList []quad.Quad
	var relationQuad quad.Quad
	var err error
	// get the relationquad
	it, _ := iterator.NewAnd(
		store,
		store.QuadIterator(quad.Object, store.ValueOf(quad.IRI(ID))),
		store.QuadIterator(quad.Predicate, store.ValueOf(model.RelationidPredicate)),
	).Optimize()
	defer it.Close()

	it, _ = store.OptimizeIterator(it)
	for it.Next() {

		relationQuad = store.Quad(it.Result())
		deleteList = append(deleteList, relationQuad)
	}
	// get the baseQuad
	baseQuad := quad.Quad{}
	for _, foundRelation := range deleteList {

		err = json.Unmarshal([]byte(foundRelation.Subject.String()), baseQuad)
		deleteList = append(deleteList, baseQuad)
	}

	if err == nil {
		// get metadata quads
		metadataQuads := GetMetadataQuadsForRelationID(ID)
		deleteList = append(deleteList, metadataQuads...)
		tx := cayley.NewTransaction()
		for _, quad := range deleteList {
			tx.RemoveQuad(quad)
		}
		err = store.ApplyTransaction(tx)
	}
	return err
}

// AddQuadRelationship will add a quad and a relationId quad to the underlying datastore
func AddQuadRelationship(relation *model.Relation) error {
	relation.ID = quad.IRI(uuid.NewUUID().String())

	relationQuad := quad.Make(quad.IRI(relation.SourceID),
		relation.Type,
		quad.IRI(relation.TargetID),
		relation.Label)

	relationSubject, parseErr := GetRelationshipSubject(relationQuad)
	if parseErr != nil {
		return parseErr
	}
	relationIDQuad := quad.Make(quad.IRI(relationSubject),
		model.RelationidPredicate,
		quad.IRI(relation.ID),
		"")

	tx := cayley.NewTransaction()

	tx.AddQuad(relationQuad)
	tx.AddQuad(relationIDQuad)

	err := store.ApplyTransaction(tx)
	return err
}
