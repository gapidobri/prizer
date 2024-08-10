package database

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)

type ParticipationMethod struct {
	Id     string              `db:"participation_method_id"`
	GameId string              `db:"game_id"`
	Name   string              `db:"name"`
	Limit  *ParticipationLimit `db:"limit"`
	Fields FieldConfig         `db:"fields"`
}

type ParticipationLimit string

const (
	ParticipationLimitDaily ParticipationLimit = "daily"
)

type FieldType string

const (
	FieldTypeString FieldType = "string"
	FieldTypeBool   FieldType = "bool"
)

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
	Type     FieldType `db:"type"`
	Required bool      `db:"required"`
	Unique   bool      `db:"unique"`
}
