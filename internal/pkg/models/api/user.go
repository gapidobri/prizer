package api

import dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"

// swagger:model User
type User struct {
	// required: true
	Id string `json:"id"`

	// required: true
	GameId string `json:"game_id"`

	// required: true
	Email *string `json:"email"`

	// required: true
	Address *string `json:"address"`

	// required: true
	Phone *string `json:"phone"`

	// required: true
	AdditionalFields map[string]any `json:"additional_fields"`
}

func UserFromDB(user dbModels.User) User {
	return User{
		Id:               user.Id,
		GameId:           user.GameId,
		Email:            user.Email,
		Address:          user.Address,
		Phone:            user.Phone,
		AdditionalFields: user.AdditionalFields,
	}
}

type GetUsersFilter struct {
	GameId *string `form:"gameId" binding:"omitnil,uuid"`
}

func (f GetUsersFilter) ToDB() dbModels.GetUsersFilter {
	return dbModels.GetUsersFilter{
		GameId: f.GameId,
	}
}

// swagger:model GetUsersResponse
type GetUsersResponse []User

// swagger:model GetUserResponse
type GetUserResponse *User
