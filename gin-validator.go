package ahlib_gin_gorm

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
func SetupGinRegexBinding() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
			return matchString(fl.Param(), fl.Field().String())
		})
	}
}

// setup binding tag: $tag
// example:
//     SetupGinSpecificRegexpBinding("phone", "^(13[0-9]|15[012356789]|17[678]|18[0-9]|14[57])[0-9]{8}$")
func SetupGinSpecificRegexpBinding(tag string, re string) {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		_ = val.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
			return matchString(re, fl.Field().String())
		})
	}
}

// setup binging tag for datetime
// example:
//     SetupGinDateTimeBinding("date", "2006-01-02", xcondition.First(time.LoadLocation("Asia/Shanghai")))
func SetupGinDateTimeBinding(tag string, layout string, loc *time.Location) {
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
