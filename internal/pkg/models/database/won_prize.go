package database

type WonPrize struct {
	Prize
	Participation
	User
}

type CreateWonPrize struct {
	PrizeId         string `db:"prize_id"`
	ParticipationId string `db:"participation_id"`
}

type GetWonPrizesFilter struct {
	GameId  *string
	UserId  *string
	PrizeId *string
}
