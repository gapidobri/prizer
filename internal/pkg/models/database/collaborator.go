package database

import "time"

type Collaborator struct {
	Id           string     `db:"collaborator_id"`
	Email        string     `db:"email"`
	Address      *string    `db:"address"`
	GameId       string     `db:"game_id"`
	LastRollTime *time.Time `db:"last_roll_time"`
}

type CreateCollaborator struct {
	GameId  string  `db:"game_id"`
	Email   string  `db:"email"`
	Address *string `db:"address"`
}
