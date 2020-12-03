package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/pquerna/ffjson/ffjson"
	"go-trailer-api/pkg/util"
)

func IntStatus(fl validator.FieldLevel) bool {
	if util.ExistIntElement(fl.Field().Int(), []int64{0, 1}) {
		return true
	}

	return false
}

func EventType(fl validator.FieldLevel) bool {
	if util.ExistIntElement(fl.Field().Int(), []int64{0, 1}) {
		return true
	}

	return false
}

func EventKt(fl validator.FieldLevel) bool {
	var (
		v   map[string]string
		err error
	)
	err = ffjson.Unmarshal([]byte(fl.Field().String()), &v)
	if err != nil {
		return false
	}

	return true
}
