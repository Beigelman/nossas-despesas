package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Validator struct {
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{validate: validator.New()}
}

func (v *Validator) Validate(data any) error {
	errs := v.validate.Struct(data)
	if errs != nil {
		var errMsgs = make([]string, len(errs.(validator.ValidationErrors)))
		for i, err := range errs.(validator.ValidationErrors) {
			errMsgs[i] = fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.Field(),
				err.Value(),
				err.Tag(),
			)
		}

		return fmt.Errorf("validation errors: %v", strings.Join(errMsgs, " and "))
	}

	return nil
}
