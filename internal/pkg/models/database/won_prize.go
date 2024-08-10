package database

type CreateWonPrize struct {
	PrizeId         string `db:"prize_id"`
	ParticipationId string `db:"participation_id"`
}
