package aceweb

import (
	"net/http"

	service "github.com/gkontos/gasket/aceservice"
	"github.com/gkontos/gasket/model"
	"github.com/gorilla/mux"
)

// RelationCreate add a relation
// return the created relation object
func RelationCreate(w http.ResponseWriter, r *http.Request) {
	var relation model.Relation
	if err := ParseJsonRequest(r, &relation); err != nil {
		ReturnErrorJSON(w, err)
		return
	}
	if err := service.AddQuadRelationship(&relation); err != nil {
		ReturnErrorJSON(w, err)
		return
	}
	ReturnBodyJSON(w, relation, http.StatusCreated)
}

// RelationDelete will delete the relation quad, and its metadata quads
func RelationDelete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var ID string
	ID = vars["id"]

	if err := service.DeleteByRelationID(ID); err != nil {
		ReturnErrorJSON(w, err)
		return
	}
	ReturnBlankJSON(w, http.StatusNoContent)
	return

}

// RelationGet will return the relation associated with the given ID
func RelationGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var ID string
	ID = vars["id"]

	q := service.GetRelation(ID)
	//	if _, ok := q.SourceID.(quad.IRI); ok { // invalid type assertion: q.SourceID.(quad.IRI) (non-interface type quad.IRI on left)
	// if q.SourceID == nil { // IRI is not type nil
	// if q.SourceID == (quad.IRI{}) { // invalid type for composite literal
	if q == (model.Relation{}) {
		ReturnBlankJSON(w, http.StatusNotFound)
		return
	}
	ReturnBodyJSON(w, q, http.StatusOK)
	return

}
