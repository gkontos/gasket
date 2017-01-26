package aceweb

import (
	"github.com/gorilla/mux"
)

var version = "v0"

func SetVersion(ver string) {
	version = ver
}

// SysViewRouter will return a router for the project routes
// TODO add auth, cors etc to handlers
// ie   router.Handle("/v1/x", common.ErrorHandler(stats.GetS)).Methods("GET")
func SysViewRouter() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(false)
	s := router.PathPrefix(version).Subrouter()

	s.HandleFunc("/nodes", NodeCreate).Methods("POST")
	s.HandleFunc("/nodes/{id}", NodeDelete).Methods("DELETE")
	s.HandleFunc("/nodes/{id}", NodeGet).Methods("GET")
	s.HandleFunc("/nodes/{id}/relationships", NodeGetRelationships).Methods("GET")
	s.HandleFunc("/nodes/{id}", NodeUpdate).Methods("PUT")

	// Given a quad, return the details of relationship
	s.HandleFunc("/relations", RelationCreate).Methods("POST")
	s.HandleFunc("/relations/{id}", RelationGet).Methods("GET")
	s.HandleFunc("/relations/{id}", RelationDelete).Methods("DELETE")
	// no PUT available for relationships.  It seems unnecessary to update a quad

	s.HandleFunc("/metadata", MetadataAdd).Methods("POST")
	// alias for /metadata endpoint
	s.HandleFunc("/relations/{id}/metadata", MetadataAdd).Methods("POST")

	s.HandleFunc("/metadata/{metadataid}", MetadataGet).Methods("GET")
	// Delete the metadata for the given quad
	s.HandleFunc("/metadata/{metadataid}", MetadataDelete).Methods("DELETE")
	// Add or update the metadata for the given quad
	s.HandleFunc("/metadata/{metadataid}", MetadataUpdate).Methods("PUT")

	return router
}
