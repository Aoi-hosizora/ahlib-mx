package xgin

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

// Setup userDefined bind.
func SetupBinding(tag string, fn func(fl validator.FieldLevel) bool) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation(tag, fn)
	}
}

// Setup binding tag: `regexp=xxx`.
func EnableRegexpBinding() {
	SetupBinding("regexp", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(fl.Param())
		return re.MatchString(fl.Field().String())
	})
}

// Setup binding tag: `$tag·.
func SetupRegexpBinding(tag string, re *regexp.Regexp) {
	SetupBinding(tag, func(fl validator.FieldLevel) bool {
		return re.MatchString(fl.Field().String())
	})
}

// Setup binging tag for datetime.
func SetupDateTimeBinding(tag string, layout string) {
	SetupBinding(tag, func(fl validator.FieldLevel) bool {
		_, err := time.Parse(layout, fl.Field().String())
		if err != nil {
			return false
		}
		return true
	})
}

// Setup length binding tag: $tag.
func SetupLengthBinding(tag string, min, max int) {
	SetupBinding(tag, func(fl validator.FieldLevel) bool {
		f := fl.Field().String()
		return len(f) >= min && len(f) <= max
	})
}

// Setup oneof binding tag: $tag.
func SetupOneofBinding(tag string, fields ...interface{}) {
	SetupBinding(tag, func(fl validator.FieldLevel) bool {
		f := fl.Field().String()
		for _, ff := range fields {
			if f == ff {
				return true
			}
		}
		return false
	})
}

// Check is err is validator.ValidationErrors and is required error.
func IsValidationFormatError(err error) bool {
	// validator.ValidationErrors
	// *errors.errorString
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	}

	for _, field := range errs {
		if field.Tag() == "required" {
			return false
		}
	}
	return true
}
