package model

// HashRequest represents the incoming JSON body for the hash endpoint.
type HashRequest struct {
	Input string `json:"input"`
}

// HashResponse represents the JSON response returned by the hash endpoint.
type HashResponse struct {
	Input string `json:"input"`
	Hash  string `json:"hash"`
}

// ErrorResponse represents a JSON error response.
type ErrorResponse struct {
	Error string `json:"error"`
}
