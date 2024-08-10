package database

type User struct {
	Id               string  `db:"user_id"`
	GameId           string  `db:"game_id"`
	AdditionalFields JsonMap `db:"additional_fields"`
	UserFields
}

type CreateUser struct {
	GameId           string  `db:"game_id"`
	AdditionalFields JsonMap `db:"additional_fields"`
	UserFields
}

type UserFields struct {
	Email   *string `db:"email"`
	Address *string `db:"address"`
	Phone   *string `db:"phone"`
}
