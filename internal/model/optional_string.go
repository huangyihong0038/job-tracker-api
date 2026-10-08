package model

import (
	"bytes"
	"encoding/json"
)

type OptionalString struct {
	Present bool
	Value   *string
}

func (field *OptionalString) UnmarshalJSON(
	data []byte,
) error {
	field.Present = true

	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		field.Value = nil
		return nil
	}

	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	field.Value = &value
	return nil

}
