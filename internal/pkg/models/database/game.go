package database

type Game struct {
	Id   string `db:"game_id"`
	Name string `db:"name"`
}
