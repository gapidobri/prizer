package enums

// swagger:enum DrawMethod
type DrawMethod string

const (
	DrawMethodFirstN DrawMethod = "first_n"
	DrawMethodChance DrawMethod = "chance"
)
