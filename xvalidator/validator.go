package xvalidator

import (
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

// Check if error is validator.ValidationErrors and with invoked by `required`.
func ValidationRequiredError(err error) bool {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// is not validation error
		return false
	}

	for _, field := range errs {
		if field.Tag() == "required" {
			// invoked by `required`
			return true
		}
	}

	// is not invoked by `required`
	return false
}

type ValidatorFunc func(fl validator.FieldLevel) bool

// ,
func And(fns ...ValidatorFunc) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		for _, fn := range fns {
			if !fn(fl) {
				return false
			}
		}
		return true
	}
}

// |
func Or(fns ...ValidatorFunc) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		for _, fn := range fns {
			if fn(fl) {
				return true
			}
		}
		return false
	}
}

func DefaultRegexpValidator() ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		param := fl.Param()
		i := fl.Field().Interface()
		str, ok := i.(string)
		if !ok {
			return false
		}

		re := regexp.MustCompile(param)
		return re.MatchString(str)
	}
}

func RegexpValidator(re *regexp.Regexp) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		str, ok := i.(string)
		if !ok {
			return false
		}
		return re.MatchString(str)
	}
}

func DateTimeValidator(tag, layout string) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		str, ok := i.(string)
		if !ok {
			return false
		}

		_, err := time.Parse(layout, str)
		if err != nil {
			return false
		}
		return true
	}
}

// Used for eq, ne.
func eqHelper(i, p interface{}) bool {
	iv, err := xreflect.IufsOf(i)
	if err != nil {
		return xreflect.IsEqual(i, p)
	}

	switch iv.Flag() {
	case xreflect.Int:
		p, ok := xreflect.GetInt(p)
		return ok && p == iv.Int()
	case xreflect.Uint:
		p, ok := xreflect.GetUint(p)
		return ok && p == iv.Uint()
	case xreflect.Float:
		p, ok := xreflect.GetFloat(p)
		return ok && xnumber.DefaultAccuracy.Equal(p, iv.Float())
	case xreflect.String:
		p, ok := xreflect.GetString(p)
		return ok && p == iv.String()
	default:
		return xreflect.IsEqual(i, p)
	}
}

// Used in len, gt, gte, lt, lte.
func lenHelper(i, p interface{}, fi func(i, p int64) bool, fu func(i, p uint64) bool, ff func(i, p float64) bool) bool {
	is, err := xreflect.IufSizeOf(i)
	if err != nil {
		return false
	}

	switch is.Flag() {
	case xreflect.Int:
		p, ok := xreflect.GetInt(p)
		return ok && fi(is.Int(), p)
	case xreflect.Uint:
		p, ok := xreflect.GetUint(p)
		return ok && fu(is.Uint(), p)
	case xreflect.Float:
		p, ok := xreflect.GetFloat(p)
		return ok && ff(is.Float(), p)
	default:
		return false
	}
}

// eq
func EqualValidator(p interface{}) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return eqHelper(i, p)
	}
}

// ne
func NotEqualValidator(p interface{}) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return !eqHelper(i, p)
	}
}

// len
func LengthValidator(p interface{}) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return lenHelper(i, p, func(i, p int64) bool {
			return i == p
		}, func(i, p uint64) bool {
			return i == p
		}, func(i, p float64) bool {
			return xnumber.DefaultAccuracy.Equal(i, p)
		})
	}
}

// gt
func GreaterThenValidator(p interface{}) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return lenHelper(i, p, func(i, p int64) bool {
			return i > p
		}, func(i, p uint64) bool {
			return i > p
		}, func(i, p float64) bool {
			return xnumber.DefaultAccuracy.Greater(i, p)
		})
	}
}

// lt
func LessThenValidator(p interface{}) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return lenHelper(i, p, func(i, p int64) bool {
			return i < p
		}, func(i, p uint64) bool {
			return i < p
		}, func(i, p float64) bool {
			return xnumber.DefaultAccuracy.Smaller(i, p)
		})
	}
}

// gte
func GreaterThenOrEqualValidator(p interface{}) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return lenHelper(i, p, func(i, p int64) bool {
			return i >= p
		}, func(i, p uint64) bool {
			return i >= p
		}, func(i, p float64) bool {
			return xnumber.DefaultAccuracy.GreaterOrEqual(i, p)
		})
	}
}

// lte
func LessThenOrEqualValidator(p interface{}) ValidatorFunc {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return lenHelper(i, p, func(i, p int64) bool {
			return i < p
		}, func(i, p uint64) bool {
			return i < p
		}, func(i, p float64) bool {
			return xnumber.DefaultAccuracy.SmallerOrEqual(i, p)
		})
	}
}
