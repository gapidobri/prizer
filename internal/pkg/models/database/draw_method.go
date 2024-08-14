package database

import "github.com/gapidobri/prizer/internal/pkg/models/enums"

type DrawMethod struct {
	Id     string           `db:"draw_method_id"`
	GameId string           `db:"game_id"`
	Name   string           `db:"name"`
	Method enums.DrawMethod `db:"method"`
	Data   string           `db:"data"`
}

type DrawMethodChanceData struct {
	Chance float64 `db:"chance"`
}

type DrawMethodFirstNData struct {
	Count int `db:"count"`
}

type GetDrawMethodsFilter struct {
	GameId                *string
	ParticipationMethodId *string
}
