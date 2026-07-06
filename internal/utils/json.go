package utils

import (
	"encoding/json"

	"github.com/vkhangstack/hexagonal-architecture/internal/core/domain"
)

func MapToStruct(m map[string]interface{}, s interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, s)
}

func ToJSONB(v any) (domain.JSONB, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return domain.JSONB(b), nil
}

func ToJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

func ParseJSONB[T any](data domain.JSONB) (T, error) {
	var out T
	err := json.Unmarshal(data, &out)
	return out, err
}

func FromJSONB[T any](data domain.JSONB, dest *T) error {
	return json.Unmarshal([]byte(data), dest)
}

func ParseJSONStringJSONB[T any](data []byte) (T, error) {
	var out T

	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return out, err
	}

	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return out, err
	}

	return out, nil
}
