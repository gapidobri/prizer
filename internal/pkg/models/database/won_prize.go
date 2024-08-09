package database

type CreateWonPrize struct {
	PrizeId        string `db:"prize_id"`
	CollaboratorId string `db:"collaborator_id"`
}
