package aceweb

import (
	"net/http"

	"fmt"

	service "github.com/gkontos/gasket/aceservice"
	"github.com/gkontos/gasket/model"
	"github.com/cayleygraph/cayley/quad"
	"github.com/gorilla/mux"
)

type RelationRequest struct {
	BaseQuad quad.Quad          `json:"base_quad"`
	Relation model.NodeProperty `json:"relation"`
}

//MetadataAdd will save a metadata struct to the store
func MetadataAdd(w http.ResponseWriter, r *http.Request) {

	var metadata model.Metadata

	parseErr := ParseJsonRequest(r, &metadata)

	if parseErr != nil {
		ReturnErrorJSON(w, parseErr)
	}

	err := service.AddMetadata(&metadata)
	if err != nil {
		ReturnErrorJSON(w, err)
		return
	}

	ReturnBodyJSON(w, metadata, http.StatusCreated)
}

func MetadataGet(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	metadataID := vars["metadataid"]

	quadList := service.GetMetadataQuadsByID(metadataID)

	if len(quadList) == 0 {
		ReturnBlankJSON(w, http.StatusNotFound)
		return
	}
	metadata, err := service.QuadListToMetadata(quadList)
	if err != nil {
		ReturnErrorJSON(w, err)
	}
	ReturnBodyJSON(w, metadata, http.StatusOK)

}

func MetadataDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	metadataID := vars["metadataid"]

	err := service.DeleteMetadataQuads(metadataID)

	if err != nil {
		ReturnErrorJSON(w, err)
		return
	}
	ReturnBlankJSON(w, http.StatusNoContent)
}

//MetadataUpdate will add or update quads associated with the metadataid.  It will return the resulting metadata struct
func MetadataUpdate(w http.ResponseWriter, r *http.Request) {
	var metadata model.Metadata

	vars := mux.Vars(r)
	metadataID := vars["metadataid"]

	parseErr := ParseJsonRequest(r, &metadata)
	if parseErr != nil {
		ReturnErrorJSON(w, parseErr)
		return
	}
	if string(metadata.ID) != "" && string(metadata.ID) != metadataID {
		validErr := &ValidationError{
			Err:     fmt.Errorf("Unable to process request"),
			Message: "Received ID's do not match",
		}
		ReturnErrorJSON(w, validErr)
		return
	} else if string(metadata.ID) == "" {
		metadata.ID = quad.IRI(metadataID)
	}

	err := service.UpdateMetadata(metadata)

	quadList := service.GetMetadataQuadsByID(string(metadata.ID))

	metadata, err = service.QuadListToMetadata(quadList)

	if err != nil {
		ReturnErrorJSON(w, err)
		return
	}
	ReturnBodyJSON(w, metadata, http.StatusOK)
	return
}
