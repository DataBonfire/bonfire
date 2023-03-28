package resource

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Filter map[string]interface{}

func (f *Filter) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var result map[string]interface{}
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		return err
	}
	for k, v := range result {
		switch t := v.(type) {
		case map[string]interface{}:
			mb, _ := json.Marshal(t)
			var c Constraint
			if err = json.Unmarshal(mb, &c); err != nil {
				return err
			}
			result[k] = &c
		}
	}

	*f = Filter(result)
	return nil
}

func (f Filter) Value() (driver.Value, error) {
	if len(f) == 0 {
		return "", nil
	}
	v, err := json.Marshal(f)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (f Filter) Match(record interface{}) bool {
	return false
}

type Constraint struct {
	Range  []interface{} `json:"range,omitempty"`
	Like   string        `json:"like,omitempty"`
	GE     interface{}   `json:"ge,omitempty"`
	GT     interface{}   `json:"gt,omitempty"`
	LE     interface{}   `json:"le,omitempty"`
	LT     interface{}   `json:"lt,omitempty"`
	Weight float32       `json:"weight,omitempty"`

	Negate bool `json:"negate,omitempty"`
}
