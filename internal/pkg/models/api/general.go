package api

// swagger:model ErrorResponse
type ErrorResponse struct {
	// required: true
	Error string `json:"error"`

	// required: true
	Code string `json:"code"`
}
