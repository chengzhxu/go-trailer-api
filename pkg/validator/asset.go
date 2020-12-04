package validator

import (
	"github.com/go-playground/validator/v10"
	"go-trailer-api/pkg/util"
	"regexp"
	"strconv"
)

var validScore = regexp.MustCompile(`^(\d|10)(\.\d)?$`)
var validPrimary = regexp.MustCompile(`^[1-9]\d*$`)

func BasPrimary(fl validator.FieldLevel) bool {
	//if validPrimary.MatchString(fl.Field().String()) {
	//	return true
	//}
	if fl.Field().Int() < 1 {
		return false
	}
	if fl.Field().Int()%1 != 0 {
		return false
	}

	return true
}

func AssetScore(fl validator.FieldLevel) bool {
	//if validPrimary.MatchString(fl.Field().String()) {
	//	return true
	//}
	if fl.Field().Float() < 0 {
		return false
	}
	if fl.Field().Float() > 10 {
		return false
	}
	str := strconv.FormatFloat(fl.Field().Float(), 'g', 5, 32)
	if len(str) > 3 && fl.Field().Float() != 10.0 {
		return false
	}

	return true
}

func AssetViewLimit(fl validator.FieldLevel) bool {
	if util.ExistIntElement(fl.Field().Int(), []int64{0, 1}) {
		return true
	}

	return false
}

func AssetType(fl validator.FieldLevel) bool {
	if util.ExistIntElement(fl.Field().Int(), []int64{1, 2}) {
		return true
	}

	return false
}

func AssetActType(fl validator.FieldLevel) bool {
	if util.ExistIntElement(fl.Field().Int(), []int64{1, 2, 3, 4}) {
		return true
	}

	return false
}
