package resource

import (
	"reflect"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(record interface{}) error {
	err := validate.Struct(record)
	if err == nil {
		return nil
	}
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}
	errs := make(map[string]string)
	t := reflect.TypeOf(record).Elem()
	for _, err := range err.(validator.ValidationErrors) {
		f, _ := t.FieldByName(err.Field())
		k := f.Tag.Get("json")
		if k == "" {
			k = err.Field()
		}
		errs[k] = err.Tag()
	}
	respErr := errors.New(400, "", "validate failure")
	respErr.Status.Metadata = errs
	return respErr
}
