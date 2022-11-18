package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Checkphrase validator.Func = func(fl validator.FieldLevel) bool {
	passphrase, ok := fl.Field().Interface().(string)
	if ok {
		words := strings.Split(passphrase, " ")
		p := asInt(fl.Param())
		return len(words) == int(p)
	}
	return false
}
