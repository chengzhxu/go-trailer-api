package validator

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Setup() {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic(errors.New("Binding Validator Engine Error\n"))
	}

	arr := map[string]func(fl validator.FieldLevel) bool{
		"bas_date":       BasDate,
		"bas_time":       BasTime,
		"nm_bas_time":    NoMustBasTime,
		"int_status":     IntStatus,
		"sdk_event_type": EventType,
		"sdk_event_kt":   EventKt,
	}

	for k, v := range arr {
		if err := validate.RegisterValidation(k, v); err != nil {
			panic(err)
		}
	}
}
