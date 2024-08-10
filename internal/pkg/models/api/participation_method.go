package api

// swagger:model ParticipationRequest
type ParticipationRequest struct {
	Fields map[string]any `json:"fields"`
}

// swagger:model ParticipationResponse
type ParticipationResponse struct {
	Prizes []PublicPrize `json:"prizes"`
}
