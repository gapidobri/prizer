package api

import (
	dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"
	"github.com/gapidobri/prizer/internal/pkg/models/enums"
)

// swagger:model DrawMethod
type DrawMethod struct {
	// required: true
	Id string `json:"draw_method_id"`

	// required: true
	GameId string `json:"game_id"`

	// required: true
	Name string `json:"name"`

	// required: true
	Method enums.DrawMethod `json:"method"`

	// required: true
	Data string `json:"data"`
}

func DrawMethodFromDB(drawMethod dbModels.DrawMethod) DrawMethod {
	return DrawMethod{
		Id:     drawMethod.Id,
		GameId: drawMethod.GameId,
		Name:   drawMethod.Name,
		Method: drawMethod.Method,
		Data:   drawMethod.Data,
	}
}

type GetDrawMethodsFilter struct {
	GameId          *string `form:"game_id" binding:"omitnil,uuid"`
	ParticipationId *string `form:"participation_id" binding:"omitnil,uuid"`
}

func (f GetDrawMethodsFilter) ToDB() dbModels.GetDrawMethodsFilter {
	return dbModels.GetDrawMethodsFilter{
		GameId:                f.GameId,
		ParticipationMethodId: f.ParticipationId,
	}
}

// swagger:model GetDrawMethodsResponse
type GetDrawMethodsResponse []DrawMethod
