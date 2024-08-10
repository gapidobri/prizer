package database

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)

type JsonMap map[string]any

func (f JsonMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &f)
}

func (f JsonMap) Value() (driver.Value, error) {
	return json.Marshal(f)
}
