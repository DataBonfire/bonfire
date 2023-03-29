package resource

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

func GormFilter(db *gorm.DB, filter *Filter) (*gorm.DB, error) {
	chains := db.Where("")
	for fieldName, fieldValue := range *filter {
		filterAssert[float64](chains, fieldValue, fmt.Sprintf("%s = ?", fieldName))
		filterAssert[string](chains, fieldValue, fmt.Sprintf("%s = ?", fieldName))
		filterAssert[bool](chains, fieldValue, fmt.Sprintf("%s = ?", fieldName))
		filterAssert[[]interface{}](chains, fieldValue, fmt.Sprintf("%s in ?", fieldName))
		if _, ok := fieldValue.(map[string]interface{}); ok {
			constraint, ok := fieldValue.(*Constraint)
			if !ok {
				continue
			}
			if constraint.Negate {
				notDb := constraintFilter(db, constraint, fieldName)
				chains.Not(notDb)
			} else {
				constraintFilter(chains, constraint, fieldName)
			}
			continue
		}

	}
	return chains, nil
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
					chains.Where(condition, constraint.Like)
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
