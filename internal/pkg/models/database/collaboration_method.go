package database

type CollaborationMethod struct {
	Id     string `db:"collaboration_method_id"`
	GameId string `db:"game_id"`
	Name   string `db:"name"`
	Fields string `db:"fields"`
}
