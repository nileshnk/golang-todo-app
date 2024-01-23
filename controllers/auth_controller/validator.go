package auth_controller

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type ValidationError struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

var Validate = validator.New()
var Trans ut.Translator

func init_validator() {
	var english = en.New()
	var uni = ut.New(english, english)
	Trans, _ = uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, Trans)
}

func translateError(err error) (errs []ValidationError) {
	if err == nil {
		return nil
	}

	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := e.Translate(Trans)
		validationError := ValidationError{
			Key:   e.Field(),
			Error: translatedErr,
		}
		errs = append(errs, validationError)
	}
	return errs
}
