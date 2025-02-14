package api

import (
	"github.com/Albertyao1993/simple_bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fiedlevel validator.FieldLevel) bool {
	if currency, ok := fiedlevel.Field().Interface().(string); ok {
		// check currency is supported
		return util.IsSupportedCurrency(currency)
	}

	return false
}
