package resource

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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

func (f Filter) Match(record interface{}, userRelation *UserRelation) bool {
	return true
	//vf := reflect.ValueOf(record)
	//if vf.Kind() == reflect.Pointer {
	//	vf = vf.Elem()
	//}
	//if vf.Kind() != reflect.Struct {
	//	return false
	//}
	//// todo 如果 record 没有 ID 或者 OrganizationID 如何处理
	//// 暂时直接返回 false
	//
	//for k, v := range f {
	//	kid, ok := reflectValueConvert(vf, k)
	//	if !ok {
	//		// todo 字段配置错误 或者 字段类型错误，需要跳过还是直接返回失败 ?
	//		return false
	//	}
	//
	//	switch v.(type) {
	//	case string:
	//		// created_by
	//		createdBy, ok := reflectValueTConvert[uint](vf, "CreatedBy")
	//		if !ok {
	//			return false
	//		}
	//		organizationID, ok := reflectValueTConvert[uint](vf, "OrganizationID")
	//		if !ok {
	//			return false
	//		}
	//		// todo 如果是其他字段, 如 campaign_id，还需要把 repo 传进来查询
	//		switch v.(string) {
	//		case "me":
	//			return createdBy == userRelation.UserId
	//		case "org":
	//			return organizationID == userRelation.OrganizationID
	//		case "sub":
	//			return isInSlice[uint](createdBy, userRelation.Subordinates)
	//		default:
	//
	//		}
	//	//case map[string]interface{}:
	//	//	mb, _ := json.Marshal(v)
	//	//	var c Constraint
	//	//	if err := json.Unmarshal(mb, &c); err != nil {
	//	//		return false
	//	//	}
	//	//case *Constraint:
	//	//	constraint := v.(*Constraint)
	//	//	if constraint.GE > 0 {
	//	//		if
	//	//	}
	//	default:
	//
	//	}
	//
	//}
	//
	//return false
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
