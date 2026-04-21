package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var currencyRe = regexp.MustCompile(`^[A-Z]{3}$`)

var Validate = newValidate()

func newValidate() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		return IsCurrency(fl.Field().String())
	})
	return v
}

func IsCurrency(s string) bool {
	return currencyRe.MatchString(s)
}
