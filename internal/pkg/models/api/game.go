package api

// swagger:model RollRequest
type RollRequest struct {
	Email   string `json:"email"`
	Address string `json:"address"`
}

// swagger:model RollResponse
type RollResponse struct {
	Won   bool   `json:"won"`
	Prize *Prize `json:"prize"`
}
