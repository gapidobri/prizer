package api

import dbModels "github.com/gapidobri/prizer/internal/pkg/models/database"

// swagger:model User
type User struct {
	Id               string         `json:"id"`
	GameId           string         `json:"game_id"`
	Email            *string        `json:"email"`
	Address          *string        `json:"address"`
	Phone            *string        `json:"phone"`
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

// swagger:model GetUsersResponse
type GetUsersResponse []User

// swagger:model GetUserResponse
type GetUserResponse *User
