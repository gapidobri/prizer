package database

import (
	"database/sql/driver"
	"encoding/json"
)

type JsonColumn[T any] struct {
	v *T
}

func (j *JsonColumn[T]) Scan(source any) error {
	if source == nil {
		j.v = nil
		return nil
	}
	j.v = new(T)
	return json.Unmarshal(source.([]byte), j.v)
}

func (j *JsonColumn[T]) Value() (driver.Value, error) {
	return json.Marshal(j.v)
}

func (j *JsonColumn[T]) Get() *T {
	return j.v
}
