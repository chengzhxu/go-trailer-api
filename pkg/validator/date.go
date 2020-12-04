package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var validBasDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
var validBasTime = regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)

func BasDate(fl validator.FieldLevel) bool {
	if validBasDate.MatchString(fl.Field().String()) {
		return true
	}

	return false
}

func BasTime(fl validator.FieldLevel) bool {
	if validBasTime.MatchString(fl.Field().String()) {
		return true
	}

	return false
}

func NoMustBasTime(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) > 0 {
		if validBasTime.MatchString(fl.Field().String()) {
			return true
		}

		return false
	}

	return true
}
