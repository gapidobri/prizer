package database

type DrawMethodEnum string

const (
	DrawMethodFirstN DrawMethodEnum = "first_n"
	DrawMethodChance DrawMethodEnum = "chance"
)

type DrawMethod struct {
	Id     string         `db:"draw_method_id"`
	Name   string         `db:"name"`
	Method DrawMethodEnum `db:"method"`
	Data   string         `db:"data"`
}

type DrawMethodChanceData struct {
	Chance float64 `db:"chance"`
}

type DrawMethodFirstNData struct {
	Count int `db:"count"`
}

type GetDrawMethodsFilter struct {
	CollaborationMethodId *string
}
