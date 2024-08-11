package api

import (
	"github.com/gapidobri/prizer/internal/pkg/models/enums"
)

// swagger:model ParticipationRequest
type ParticipationRequest struct {
	// required: true
	Fields map[string]any `json:"fields"`
}

// swagger:model ParticipationResponse
type ParticipationResponse struct {
	// required: true
	Prizes []PublicPrize `json:"prizes"`
}

// swagger:model ParticipationMethod
type ParticipationMethod struct {
	// required: true
	Id string `json:"id"`

	// required: true
	GameId string `json:"game_id"`

	// required: true
	Name string `json:"name"`

	// required: true
	Limit *enums.ParticipationLimit `json:"limit"`

	// required: true
	Fields FieldConfig `json:"fields"`
}

type FieldConfig struct {
	// required: true
	User map[string]Field `json:"user"`

	// required: true
	Participation map[string]Field `json:"participation"`
}

type Field struct {
	// required: true
	Type enums.FieldType `json:"type"`

	// required: true
	Required bool `json:"required"`

	// required: true
	Unique bool `json:"unique"`
}

type GetParticipationMethodsFilter struct {
	GameId *string `form:"gameId" binding:"omitnil,uuid"`
}

// swagger:model GetParticipationMethodsResponse
type GetParticipationMethodsResponse []ParticipationMethod
