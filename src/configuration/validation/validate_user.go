package validation

import (
	"crud-go/src/configuration/rest_err"
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_Translation "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate = validator.New()
	transl   ut.Translator
)

func init() {
	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		un := ut.New(en, en)
		transl, _ = un.GetTranslator("en")
		en_Translation.RegisterDefaultTranslations(val, transl)
	}
}

func ValidateUserError(valition_err error) *rest_err.RestErr {
	var jsonErr *json.UnmarshalTypeError
	var jsonValidationError validator.ValidationErrors

	if errors.As(valition_err, &jsonErr) {
		return rest_err.NewBadRequestError("Invalid field type")
	} else if errors.As(valition_err, &jsonValidationError) {
		errorsCauses := []rest_err.Causes{}

		for _, e := range valition_err.(validator.ValidationErrors) {
			cause := rest_err.Causes{
				Message: e.Translate(transl),
				Field:   e.Field(),
			}

			errorsCauses = append(errorsCauses, cause)
		}

		return rest_err.NewBadRequestValidationError("Some fields are invalid", errorsCauses)
	} else {
		return rest_err.NewBadRequestError("Error trying to covert fields")
	}
}
