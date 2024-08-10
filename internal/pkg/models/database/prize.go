package database

type Prize struct {
	Id          string `db:"prize_id"`
	GameId      string `db:"game_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Count       int    `db:"count"`
}

type GetPrizesFilter struct {
	GameId                *string
	ParticipationMethodId *string
	DrawMethodId          *string
	AvailableOnly         bool
}
