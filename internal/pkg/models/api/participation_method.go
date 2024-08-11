package api

import (
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/gapidobri/prizer/internal/pkg/models/enums"
	"github.com/samber/lo"
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
	Limit enums.ParticipationLimit `json:"limit"`

	// required: true
	Fields FieldConfig `json:"fields"`
}

func ParticipationMethodFromDB(method dbModels.ParticipationMethod) ParticipationMethod {
	return ParticipationMethod{
		Id:     method.Id,
		GameId: method.GameId,
		Name:   method.Name,
		Limit:  method.Limit,
		Fields: FieldConfigFromDB(method.Fields),
	}
}

type FieldConfig struct {
	// required: true
	User map[string]Field `json:"user"`

	// required: true
	Participation map[string]Field `json:"participation"`
}

func FieldConfigFromDB(config dbModels.FieldConfig) FieldConfig {
	return FieldConfig{
		User: lo.MapValues(config.User, func(field dbModels.Field, key string) Field {
			return FieldFromDB(field)
		}),
		Participation: lo.MapValues(config.Participation, func(field dbModels.Field, key string) Field {
			return FieldFromDB(field)
		}),
	}
}

type Field struct {
	// required: true
	Type enums.FieldType `json:"type"`

	// required: true
	Required bool `json:"required"`

	// required: true
	Unique bool `json:"unique"`
}

func FieldFromDB(field dbModels.Field) Field {
	return Field{
		Type:     field.Type,
		Required: field.Required,
		Unique:   field.Unique,
	}
}

type GetParticipationMethodsFilter struct {
	GameId *string `form:"gameId" binding:"omitnil,uuid"`
}

func (f GetParticipationMethodsFilter) ToDB() dbModels.GetParticipationMethodsFilter {
	return dbModels.GetParticipationMethodsFilter{
		GameId: f.GameId,
	}
}

// swagger:model GetParticipationMethodsResponse
type GetParticipationMethodsResponse []ParticipationMethod
