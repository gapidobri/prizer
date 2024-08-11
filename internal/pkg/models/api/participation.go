package api

import (
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"time"
)

// swagger:model Participation
type Participation struct {
	// required: true
	ParticipationId string `json:"participation_id"`

	// required: true
	ParticipationMethodId string `json:"participation_method_id"`

	// required: true
	UserId string `json:"user_id"`

	// required: true
	CreatedAt time.Time `json:"created_at"`

	// required: true
	Fields map[string]any `json:"fields"`
}

func ParticipationFromDB(participation dbModels.Participation) Participation {
	return Participation{
		ParticipationId:       participation.Id,
		ParticipationMethodId: participation.Id,
		UserId:                participation.UserId,
		CreatedAt:             participation.CreatedAt,
		Fields:                participation.Fields,
	}
}
