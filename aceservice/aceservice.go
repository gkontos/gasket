package aceservice

import "github.com/cayleygraph/cayley"

var store *cayley.Handle

func SetStore(s *cayley.Handle) {
	store = s
}

type DataStoreError struct {
	Err     error
	Message string
}

func (e *DataStoreError) Error() string {
	if e.Message != "" {
		return e.Message + " " + e.Err.Error()
	}
	return e.Err.Error()
}
