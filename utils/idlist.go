package utils

import (
	"encoding/json"
	"fmt"
)

// ParseInt64ListJSON parses a JSON array of IDs (string or number) into a slice of int64.
// Returns nil on empty or null.
func ParseInt64ListJSON(data []byte, field string) ([]int64, error) {
	if len(data) == 0 || string(data) == "null" {
		return nil, nil
	}

	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%s must be an array of IDs", field)
	}

	ids := make([]int64, 0, len(raw))
	for i, item := range raw {
		var str string
		if err := json.Unmarshal(item, &str); err == nil {
			id, err := parseStringID(str)
			if err != nil {
				return nil, fmt.Errorf("%s[%d] must be a valid integer string", field, i)
			}
			ids = append(ids, id)
			continue
		}

		var num json.Number
		if err := json.Unmarshal(item, &num); err == nil {
			id, err := num.Int64()
			if err != nil || id <= 0 {
				return nil, fmt.Errorf("%s[%d] must be a positive integer", field, i)
			}
			ids = append(ids, id)
			continue
		}

		return nil, fmt.Errorf("%s[%d] must be a string or number", field, i)
	}

	return ids, nil
}

func parseStringID(s string) (int64, error) {
	id, err := json.Number(s).Int64()
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid id")
	}
	return id, nil
}

// IDList is a helper slice that can unmarshal JSON arrays of ints/strings into int64s.
type IDList []int64

// UnmarshalJSON implements custom unmarshalling with validation.
func (l *IDList) UnmarshalJSON(data []byte) error {
	ids, err := ParseInt64ListJSON(data, "ids")
	if err != nil {
		return err
	}
	*l = ids
	return nil
}

// Int64s returns a plain slice copy.
func (l IDList) Int64s() []int64 {
	return []int64(l)
}
