package aceweb

//RequestParseError is an error type indicating that there was a problem parsing the http request
type RequestParseError struct {
	Err     error
	Message string
}

func (e *RequestParseError) Error() string {
	if e.Message != "" {
		return e.Message + " " + e.Err.Error()
	}
	return e.Err.Error()
}
