package api

// swagger:model ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
