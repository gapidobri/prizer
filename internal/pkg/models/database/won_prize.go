package database

type CreateWonPrize struct {
	PrizeId string `db:"prize_id"`
	UserId  string `db:"user_id"`
}
