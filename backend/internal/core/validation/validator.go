package validation

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/jethroTLOU/tokigura/internal/core/apperr"
)

type Validator struct {
	v     *validator.Validate
	trans ut.Translator
}

func New() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(v, trans)

	return &Validator{
		v:     v,
		trans: trans,
	}
}

func (v *Validator) Struct(obj any) error {
	err := v.v.Struct(obj)
	if err == nil {
		return nil
	}

	var valErrs validator.ValidationErrors
	if !errors.As(err, &valErrs) {
		return apperr.Wrap(err, apperr.CodeInternal, "Given non validation error type")
	}

	errs := make(map[string]any, len(valErrs))
	for _, fieldError := range valErrs {
		errs[fieldError.Field()] = fieldError.Translate(v.trans)
	}

	return apperr.NewValidation(errs)
}
