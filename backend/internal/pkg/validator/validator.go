package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
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
		var errMsgs = make([]string, len(func() validator.ValidationErrors {
			var target validator.ValidationErrors
			_ = errors.As(errs, &target)
			return target
		}()))
		for i, err := range func() validator.ValidationErrors {
			var target validator.ValidationErrors
			_ = errors.As(errs, &target)
			return target
		}() {
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
