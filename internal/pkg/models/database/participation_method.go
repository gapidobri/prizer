package database

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/gapidobri/prizer/internal/pkg/models/enums"
	"github.com/pkg/errors"
)

type ParticipationMethod struct {
	Id                          string                   `db:"participation_method_id"`
	GameId                      string                   `db:"game_id"`
	Name                        string                   `db:"name"`
	ParticipationLimit          enums.ParticipationLimit `db:"participation_limit"`
	Fields                      FieldConfig              `db:"fields"`
	WinMailTemplateId           *string                  `db:"win_mail_template_id"`
	LoseMailTemplateId          *string                  `db:"lose_mail_template_id"`
	ParticipationMailTemplateId *string                  `db:"participation_mail_template_id"`
}

type FieldConfig struct {
	User          map[string]Field `json:"user"`
	Participation map[string]Field `json:"participation"`
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
	Type         enums.FieldType `json:"type"`
	Required     bool            `json:"required"`
	Unique       bool            `json:"unique"`
	MailVariable *string         `json:"mail_variable"`
}

type GetParticipationMethodsFilter struct {
	GameId *string
}

type UpdateParticipationMethod struct {
	Name                        string                   `mapstructure:"name"`
	ParticipationLimit          enums.ParticipationLimit `mapstructure:"participation_limit"`
	Fields                      FieldConfig              `mapstructure:"-"`
	WinMailTemplateId           *string                  `mapstructure:"win_mail_template_id"`
	LoseMailTemplateId          *string                  `mapstructure:"lose_mail_template_id"`
	ParticipationMailTemplateId *string                  `mapstructure:"participation_mail_template_id"`
}
