package database

type Game struct {
	Id                     string  `db:"game_id"`
	Name                   string  `db:"name"`
	WinPercentage          float32 `db:"win_percentage"`
	UniqueCollaboratorData bool    `db:"unique_collaborator_data"`
}
