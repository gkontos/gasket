package aceweb

import (
	"net/http"

	"fmt"

	log "github.com/gkontos/gasket/acelog"
	service "github.com/gkontos/gasket/aceservice"
	"github.com/gkontos/gasket/model"
	"github.com/cayleygraph/cayley/quad"
	"github.com/gorilla/mux"
)

// NodeCreate expects to receive a json node struct.  The node will be added to the store
func NodeCreate(w http.ResponseWriter, r *http.Request) {

	var node model.Node

	parseErr := ParseJsonRequest(r, &node)

	if parseErr != nil {
		log.Error(parseErr)
		ReturnErrorJSON(w, parseErr)
		return
	}

	// to generate the node id, get the type.
	// foreach other property create a quad with node id as subject and the NodeProperties as the remaining quad values
	// call AddQuad to save the properties

	quadList, err := service.AddNode(node)

	if err != nil {
		ReturnErrorJSON(w, err)
		return
	}
	node, err = service.QuadListToNode(quadList)
	if err != nil {
		ReturnErrorJSON(w, err)
		return
	}
	ReturnBodyJSON(w, node, http.StatusCreated)
	return

}

// NodeDelete will delete all quads for the object node specified by the {id}
// router.HandleFunc("/nodes/{id}", NodeDelete).Methods("DELETE")
func NodeDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var nodeID string
	nodeID = vars["id"]

	deleteErr := service.DeleteByID(nodeID)
	if deleteErr != nil {
		ReturnErrorJSON(w, deleteErr)
		return
	}

	ReturnBlankJSON(w, http.StatusNoContent)

}

// NodeGet will get the quads relating to the node specified by the {id}
// router.HandleFunc("/nodes/{id}", NodeGet).Methods("GET")
func NodeGet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var subject string
	subject = vars["id"]

	quadList := service.GetQuadsBySubject(subject)

	if len(quadList) == 0 {
		ReturnBlankJSON(w, http.StatusNotFound)
		return
	}
	node, err := service.QuadListToNode(quadList)
	if err != nil {
		ReturnErrorJSON(w, err)
	} else {
		ReturnBodyJSON(w, node, http.StatusOK)
	}
	return
}

// NodeGetRelationships will get the quads relating to the node specified by the {id}
// router.HandleFunc("/nodes/{id}/relationships", NodeGet).Methods("GET")
func NodeGetRelationships(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var subject string
	subject = vars["id"]

	quadList := service.GetQuads(subject)
	if len(quadList) == 0 {
		ReturnBlankJSON(w, http.StatusNoContent)
		return
	}
	ReturnBodyJSON(w, quadList, http.StatusOK)
	return
}

// NodeUpdate will update or add quad properties for the {id}
// router.HandleFunc("/nodes/{id}", NodeUpdateProperty).Methods("PUT")
func NodeUpdate(w http.ResponseWriter, r *http.Request) {

	var node model.Node

	vars := mux.Vars(r)
	nodeID := vars["id"]

	parseErr := ParseJsonRequest(r, &node)

	if parseErr != nil {
		ReturnErrorJSON(w, parseErr)
		return
	}

	if string(node.ID) != "" && string(node.ID) != nodeID {

		validErr := &ValidationError{
			Err:     fmt.Errorf("Unable to process request"),
			Message: "Received ID's do not match",
		}
		ReturnErrorJSON(w, validErr)
		return
	} else if string(node.ID) == "" {

		node.ID = quad.IRI(nodeID)
	}
	_, err := service.UpdateNode(node)
	if err != nil {
		ReturnErrorJSON(w, err)
		return
	}
	quadList := service.GetQuadsBySubject(nodeID)

	node, mappingErr := service.QuadListToNode(quadList)
	if mappingErr != nil {
		ReturnErrorJSON(w, mappingErr)
		return
	} else {
		ReturnBodyJSON(w, node, http.StatusOK)
	}

	return
}
