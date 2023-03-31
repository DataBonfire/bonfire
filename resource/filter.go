package resource

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm/schema"
	"reflect"
	"strconv"
)

type Filter map[string]interface{}

func (f *Filter) UnmarshalJSON(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
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

func (f *Filter) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return f.UnmarshalJSON(bytes)
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

// 1. me 2. org 3. sub (Subordinate)
// 不考虑同事

// {"action:"edit", "resource": "campaigns", "record": {"created_by.manager_id": "me", "status": {"lt": 3}}
// {"action:"edit", "resource": "campaigns", "record": {"created_by": "me"}, "status": {"lt": 3}}
// {"action:"edit", "resource": "campaigns", "record": {"created_by": "me"}, "created_by": {"lt": 3}}
// {"action:"edit", "resource": "campaigns", "record": {"created_by": "org"}, "status": {"lt": 3}}
//{"action:"edit", "resource": "campaigns", "record": {"campaign_id": "mine"}, "status": {"lt": 3}}
// campaigns -> list -> show edit btn?
// hasPermission(action, resource)

type UserRelation struct {
	UserId         uint
	OrganizationID uint
	Subordinates   []uint
}

func (f Filter) Match(record interface{}) bool {

	recordeReflectType := reflect.TypeOf(record)

	if recordeReflectType.Kind() == reflect.Pointer {
		recordeReflectType = recordeReflectType.Elem()
	}
	// record 不是 reflect.Struct ，返回 false
	if recordeReflectType.Kind() != reflect.Struct {
		return false
	}
	for i, n := 0, recordeReflectType.NumField(); i < n; i++ {
		structFieldName := recordeReflectType.Field(i).Name
		fieldName := schema.NamingStrategy{}.ColumnName("", structFieldName)
		fv, ok := f[fieldName]
		if !ok {
			continue
		}
		vf := reflect.ValueOf(record)
		if vf.Kind() == reflect.Pointer {
			vf = vf.Elem()
		}

		if vf.Kind() != reflect.Struct {
			continue
		}
		// 目前只处理 uint
		// 如 influencer_id 1
		id, ok := reflectValueTConvert[uint](vf, structFieldName)
		if !ok {
			continue
		}
		rv := reflect.TypeOf(fv)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			if isInSlice[uint](id, fv.([]uint)) {
				return true
			}
		case reflect.Uint:
			if id == fv {
				return true
			}
		}
	}
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

func reflectValueTConvert[T any](vf reflect.Value, name string) (T, bool) {
	var data T
	fieldByName := vf.FieldByName(name)
	if !fieldByName.IsValid() || !fieldByName.CanInterface() {
		return data, false
	}
	data, ok := fieldByName.Interface().(T)
	return data, ok
}

func reflectValueConvert(vf reflect.Value, name string) (interface{}, bool) {
	v := vf.FieldByName(name)
	if !v.IsValid() || !v.CanInterface() {
		return nil, false
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return v.Interface().(float64), true
	case reflect.String:
		d, err := strconv.Atoi(v.Interface().(string))
		if err != nil {
			return nil, false
		}
		return uint(d), true
	//case reflect.Bool:
	//	// xx

	default:
		return nil, false
	}
}

func isInSlice[T comparable](element T, arr []T) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}
