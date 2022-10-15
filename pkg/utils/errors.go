package utils

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func TranslateError(err error, ts ut.Translator) error {
	if err == nil {
		return nil
	}
	var finalErr string
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		finalErr += fmt.Sprintf("%s", e.Translate(ts))
	}
	return fmt.Errorf("%s", finalErr)
}
