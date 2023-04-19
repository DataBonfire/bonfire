package filter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

func GormOrder(db *gorm.DB, field, order string) (*gorm.DB, error) {
	order = strings.ToUpper(order)
	if order != "DESC" && order != "ASC" {
		return nil, errors.New("order error")
	}

	return db.Order(fmt.Sprintf("%s %v", field, order)), nil
}

func GormFilter(db *gorm.DB, filters ...Filter) (*gorm.DB, error) {
	groupDB := db.Where("")
	for _, filter := range filters {
		chains := db.Where("")
		for fieldName, v := range filter {
			rv := reflect.TypeOf(v)
			switch rv.Kind() {
			case reflect.Slice, reflect.Array:
				chains.Where(fmt.Sprintf("%s in ?", fieldName), v)
			case reflect.Pointer:
				constraint, ok := v.(*Constraint)
				if !ok {
					continue
				}
				if constraint.Negate {
					notDb := constraintFilter(db, constraint, fieldName)
					chains.Not(notDb)
				} else {
					constraintFilter(chains, constraint, fieldName)
				}
			default:
				if v == -1 && fieldName == "created_by" {
					chains.Where("FALSE")
				} else {
					chains.Where(fmt.Sprintf("%s = ?", fieldName), v)
				}
			}
		}
		groupDB.Or(chains)
	}
	return groupDB, nil
}

func filterAssert[T any](chains *gorm.DB, fieldValue interface{}, condition string) {
	if _, ok := fieldValue.(T); ok {
		chains.Where(condition, fieldValue)
	}
}

func constraintFilter(chains *gorm.DB, constraint *Constraint, fieldName string) *gorm.DB {
	vf := reflect.TypeOf(*constraint)
	if vf.Kind() == reflect.Struct {
		for i, n := 0, vf.NumField(); i < n; i++ {
			opName := vf.Field(i).Name
			switch opName {
			case "LE":
				filterCondition(chains, constraint.LE, fieldName, "%s <= ?")
			case "LT":
				filterCondition(chains, constraint.LT, fieldName, "%s < ?")
			case "GE":
				filterCondition(chains, constraint.GE, fieldName, "%s >= ?")
			case "GT":
				filterCondition(chains, constraint.GT, fieldName, "%s > ?")
			case "Like":
				if len(constraint.Like) > 0 {
					condition := fmt.Sprintf("%s like ?", fieldName)
					chains.Where(condition, "%"+constraint.Like+"%")
				}
			case "Range":
				if len(constraint.Range) == 2 {
					condition := fmt.Sprintf("%s BETWEEN ? AND ?", fieldName)
					chains.Where(condition, constraint.Range...)
				}
			}
		}
	}
	return chains
}

func filterCondition(chains *gorm.DB, value interface{}, filedName string, conditionFormat string) {
	if value == nil {
		return
	}
	switch rv := reflect.TypeOf(value); rv.Kind() {
	case reflect.Float32, reflect.Float64:
		if value.(float64) <= 0 {
			return
		}
	case reflect.String:
		if len(value.(string)) == 0 {
			return
		}
	}
	chains.Where(fmt.Sprintf(conditionFormat, filedName), value)
}
