package aceservice

import (
	"reflect"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph/iterator"
	"github.com/cayleygraph/cayley/quad"
	log "github.com/gkontos/gasket/acelog"
	"github.com/gkontos/gasket/model"
	"github.com/pborman/uuid"
)

// AddMetadata will save a metadata relation quad and the property quads to the underlying store
func AddMetadata(metadata *model.Metadata) error {
	metadata.ID = quad.IRI(uuid.NewUUID().String())
	metadataIDQuad := quad.Make(metadata.RelationID,
		model.MetaidPredicate,
		metadata.ID,
		"")

	metadataQuads := getMetadataPropertiesAsQuads(*metadata)
	metadataQuads = append(metadataQuads, metadataIDQuad)

	tx := cayley.NewTransaction()
	for _, q := range metadataQuads {
		tx.AddQuad(q)
	}
	err := store.ApplyTransaction(tx)
	return err
}

// getMetadataPropertiesAsQuads will return the properties map as a list of quads
func getMetadataPropertiesAsQuads(metadata model.Metadata) []quad.Quad {
	var metadataQuads []quad.Quad
	for key, value := range metadata.Properties {
		objectValue, ok := quad.AsValue(value)
		if !ok {
			log.Error("Unable to parse value : ", value)
		}
		q := quad.Make(
			metadata.ID,
			quad.IRI(key),
			objectValue,
			"",
		)
		metadataQuads = append(metadataQuads, q)
	}

	return metadataQuads
}

// DeleteMetadataQuads will delete the metadata relation quad and any property quads for the given Id
func DeleteMetadataQuads(metadataID string) error {
	quadList := GetMetadataQuadsByID(metadataID)
	tx := cayley.NewTransaction()
	for _, q := range quadList {
		tx.RemoveQuad(q)
	}
	err := store.ApplyTransaction(tx)
	return err
}

// UpdateMetadata will add or update any properties of the metadata object
func UpdateMetadata(metadata model.Metadata) error {
	quadList := getMetadataPropertiesAsQuads(metadata)

	tx := cayley.NewTransaction()
	for _, q := range quadList {
		AddOrUpdateAsTransaction(tx, q)
	}
	err := store.ApplyTransaction(tx)
	if err != nil {
		return &DataStoreError{Message: "Error saving data", Err: err}
	}
	return nil
}

// GetMetadataQuadsByID gets the identity quad as well as property quads for a given metadataId
func GetMetadataQuadsByID(metadataID string) []quad.Quad {
	var metaQuadList []quad.Quad
	it, _ := iterator.NewAnd(
		store,
		store.QuadIterator(quad.Object, store.ValueOf(quad.IRI(metadataID))),
		store.QuadIterator(quad.Predicate, store.ValueOf(model.MetaidPredicate)),
	).Optimize()
	defer it.Close()

	it, _ = store.OptimizeIterator(it)

	for it.Next() {

		metaQuad := store.Quad(it.Result())
		metaQuadList = append(metaQuadList, metaQuad)
	}

	for _, direction := range []quad.Direction{quad.Subject} {
		metait := store.QuadIterator(direction, store.ValueOf(quad.IRI(metadataID)))
		defer metait.Close()
		for metait.Next() {
			meta := store.Quad(metait.Result())
			log.Debug("type found ", reflect.TypeOf(meta.Object))
			log.Debug("has value ", meta.Object.String())
			metaQuadList = append(metaQuadList, meta)
		}
	}

	return metaQuadList
}

// GetMetadataQuadsForRelationId will return all metadata quads and the metadata relations for a given relationId
func GetMetadataQuadsForRelationID(relationID string) []quad.Quad {
	var metaQuadList []quad.Quad
	var metaIdList []quad.Value
	it, _ := iterator.NewAnd(
		store,
		store.QuadIterator(quad.Subject, store.ValueOf(quad.IRI(relationID))),
		store.QuadIterator(quad.Predicate, store.ValueOf(model.MetaidPredicate)),
	).Optimize()
	defer it.Close()

	it, _ = store.OptimizeIterator(it)

	for it.Next() {

		metaQuad := store.Quad(it.Result())
		metaIdList = append(metaIdList, metaQuad.Subject)
	}
	for _, metaid := range metaIdList {
		for _, direction := range []quad.Direction{quad.Subject} {
			metait := store.QuadIterator(direction, store.ValueOf(metaid))
			defer metait.Close()
			for metait.Next() {
				metaQuadList = append(metaQuadList, store.Quad(it.Result()))
			}
		}
	}
	return metaQuadList
}
