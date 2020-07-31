package xgin

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
	"time"
)

func matchString(reg string, content string) bool {
	re, err := regexp.Compile(reg)
	if err != nil {
		return true // error reg default to match success
	}
	return re.MatchString(content)
}

// setup binding tag: regexp=xxx
// noinspection GoUnusedExportedFunction
func SetupRegexBinding() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
			return matchString(fl.Param(), fl.Field().String())
		})
	}
}

// setup binding tag: $tag
// example:
//     SetupSpecificRegexpBinding("phone", "^(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$")
// noinspection GoUnusedExportedFunction
func SetupSpecificRegexpBinding(tag string, re string) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
			return matchString(re, fl.Field().String())
		})
	}
}

// setup binging tag for datetime with loc
// noinspection GoUnusedExportedFunction
func SetupDateTimeLocBinding(tag string, layout string, loc *time.Location) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
			_, err := time.ParseInLocation(layout, fl.Field().String(), loc)
			if err != nil {
				return false
			}
			return true
		})
	}
}

// setup binging tag for datetime
// noinspection GoUnusedExportedFunction
func SetupDateTimeBinding(tag string, layout string) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
			_, err := time.Parse(layout, fl.Field().String())
			if err != nil {
				return false
			}
			return true
		})
	}
}

// setup userDefined bind
// example:
//     SetupBinding("xxx", func(fl validator.FieldLevel) {
//         return fl.Field.String() == "xxx"
//     })
// noinspection GoUnusedExportedFunction
func SetupBinding(tag string, valFunc func(fl validator.FieldLevel) bool) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation(tag, valFunc)
	}
}

// noinspection GoUnusedExportedFunction
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
