package resource

import (
	"reflect"

	"github.com/go-kratos/kratos/v2/errors"
	bizvalidator "github.com/go-playground/validator/v10"
)

var validate = bizvalidator.New()

func Validate(record interface{}) error {
	err := validate.Struct(record)
	if err == nil {
		return nil
	}
	if _, ok := err.(*bizvalidator.InvalidValidationError); ok {
		return err
	}
	errs := make(map[string]string)
	t := reflect.TypeOf(record).Elem()
	for _, err := range err.(bizvalidator.ValidationErrors) {
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
