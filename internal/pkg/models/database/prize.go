package database

type Prize struct {
	Id          string  `db:"prize_id"`
	GameId      string  `db:"game_id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	ImageUrl    *string `db:"image_url"`
	Count       int     `db:"count"`
	WonCount    int     `db:"won_count"`
}

type GetPrizesFilter struct {
	GameId        *string
	DrawMethodId  *string
	UserId        *string
	AvailableOnly bool
}

type CreatePrize struct {
	GameId      string  `db:"game_id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	ImageUrl    *string `db:"image_url"`
	Count       int     `db:"count"`
}
