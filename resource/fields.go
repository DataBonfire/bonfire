package resource

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	JSON        json.RawMessage
	StringSlice []string
	Timestamp   int64
)

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

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

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	if *t == 0 {
		return []byte("\"\""), nil
	}
	return []byte(strconv.Quote(time.Unix(int64(*t), 0).Format("2006-01-02T15:04:05.000Z"))), nil
	//return []byte(strconv.Quote(time.Unix(int64(*t), 0).Format("2006-01-02 15:04:05"))), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	v, err := strconv.Unquote(string(data))
	if err != nil || v == "" {
		*t = 0
		return nil
	}
	// 2016-01-02
	if len(v) < 10 {
		return ErrInvalidTimeFormat
	}
	//v = (strings.ReplaceAll(v, "/", "-") + " 00:00:00")[:19]
	parsed, err := time.Parse("2006-01-02T15:04:05.000Z", v)
	if err != nil {
		return err
	}
	*t = Timestamp(parsed.Unix())
	return nil
}

var ErrInvalidTimeFormat = errors.New("invalid time format")
