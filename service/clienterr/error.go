package clienterr

// Used to bubble lower-level errors up to HTTP handler
type ClientError struct {
	e      string
	Parent error