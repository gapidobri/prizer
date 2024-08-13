package database

type Game struct {
	Id                 string  `db:"game_id"`
	Name               string  `db:"name"`
	GoogleSheetId      *string `db:"google_sheet_id"`
	GoogleSheetTabName *string `db:"google_sheet_tab_name"`
}
