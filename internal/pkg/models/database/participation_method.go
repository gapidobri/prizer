package database

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/gapidobri/prizer/internal/pkg/models/enums"
	"github.com/pkg/errors"
)

type ParticipationMethod struct {
	Id     string                   `db:"participation_method_id"`
	GameId string                   `db:"game_id"`
	Name   string                   `db:"name"`
	Limit  enums.ParticipationLimit `db:"limit"`
	Fields FieldConfig              `db:"fields"`
}

type FieldConfig struct {
	User          map[string]Field `db:"user"`
	Participation map[string]Field `db:"participation"`
}

func (f *FieldConfig) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &f)
}

func (f *FieldConfig) Value() (driver.Value, error) {
	return json.Marshal(f)
}

type Field struct {
	Type     enums.FieldType `db:"type"`
	Required bool            `db:"required"`
	Unique   bool            `db:"unique"`
}

type GetParticipationMethodsFilter struct {
	GameId *string
}
