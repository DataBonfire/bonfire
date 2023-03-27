package resource

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

type StringSlice []string

func (s *StringSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to parse StringSlice value:", value))
	}

	if v := strings.TrimSpace(string(bytes)); v != "" {
		*s = strings.Split(v, ",")
	}
	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "", nil
	}
	return strings.Join(s, ","), nil
}
