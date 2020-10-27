package xvalidator

import (
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/go-playground/validator/v10"
	"reflect"
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

// ,
func And(fns ...validator.Func) validator.Func {
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
func Or(fns ...validator.Func) validator.Func {
	return func(fl validator.FieldLevel) bool {
		for _, fn := range fns {
			if fn(fl) {
				return true
			}
		}
		return false
	}
}

// regexp: xxx
func DefaultRegexpValidator() validator.Func {
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

func RegexpValidator(re *regexp.Regexp) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		str, ok := i.(string)
		if !ok {
			return false
		}
		return re.MatchString(str)
	}
}

func DateTimeValidator(layout string) validator.Func {
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
// For numbers & strings, it validates the value.
// For slices, arrays, and maps, it validates the length.
func eqHelper(i, p interface{}) bool {
	iv, err := xreflect.IufsOf(i)
	if err != nil {
		val := reflect.ValueOf(i)
		if knd := val.Kind(); knd == reflect.Slice || knd == reflect.Array || knd == reflect.Map { // slice array map
			p, ok := xreflect.GetInt(p)
			if ok {
				return val.Len() == int(p)
			}
		}
		return false
	}

	switch iv.Flag() {
	case xreflect.Int: // int
		p, ok := xreflect.GetInt(p)
		return ok && p == iv.Int()
	case xreflect.Uint: // uint
		p, ok := xreflect.GetUint(p)
		return ok && p == iv.Uint()
	case xreflect.Float: // float
		p, ok := xreflect.GetFloat(p)
		return ok && xnumber.DefaultAccuracy.Equal(p, iv.Float())
	case xreflect.String: // string
		p, ok := xreflect.GetString(p)
		return ok && p == iv.String()
	default:
		return false
	}
}

// Used in len, gt, gte, lt, lte.
// For numbers, it validates the value.
// For strings, it validates the length of string.
// For slices, arrays, and maps, it validates the length.
func lenHelper(i, p interface{}, fi func(i, p int64) bool, fu func(i, p uint64) bool, ff func(i, p float64) bool) bool {
	is, err := xreflect.IufSizeOf(i)
	if err != nil {
		return false
	}

	switch is.Flag() {
	case xreflect.Int: // int, string, slice, array, map
		p, ok := xreflect.GetInt(p)
		return ok && fi(is.Int(), p)
	case xreflect.Uint: // uint
		p, ok := xreflect.GetUint(p)
		return ok && fu(is.Uint(), p)
	case xreflect.Float: // float
		p, ok := xreflect.GetFloat(p)
		return ok && ff(is.Float(), p)
	default:
		return false
	}
}

// eq. See https://godoc.org/github.com/go-playground/validator#hdr-Equals.
// For strings & numbers, eq will ensure that the value is equal to the parameter given.
// For slices, arrays, and maps, validates the number of items.
func EqualValidator(p interface{}) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return eqHelper(i, p)
	}
}

// ne. See https://godoc.org/github.com/go-playground/validator#hdr-Not_Equal.
// For strings & numbers, ne will ensure that the value is not equal to the parameter given.
// For slices, arrays, and maps, validates the number of items.
func NotEqualValidator(p interface{}) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return !eqHelper(i, p)
	}
}

// len. See https://godoc.org/github.com/go-playground/validator#hdr-Length.
// For numbers, length will ensure that the value is equal to the parameter given.
// For strings, it checks that the string length is exactly that number of characters.
// For slices, arrays, and maps, validates the number of items.
func LenValidator(p interface{}) validator.Func {
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

// gt. See https://godoc.org/github.com/go-playground/validator#hdr-Greater_Than.
// For numbers, this will ensure that the value is greater than the parameter given.
// For strings, it checks that the string length is greater than that number of characters.
// For slices, arrays and maps it validates the number of items.
func GreaterThenValidator(p interface{}) validator.Func {
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

// lt. See https://godoc.org/github.com/go-playground/validator#hdr-Less_Than.
// For numbers, this will ensure that the value is less than the parameter given.
// For strings, it checks that the string length is less than that number of characters.
// For slices, arrays, and maps it validates the number of items.
func LessThenValidator(p interface{}) validator.Func {
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

// gte. See https://godoc.org/github.com/go-playground/validator#hdr-Greater_Than_or_Equal.
// For numbers, gte will ensure that the value is greater or equal to the parameter given.
// For strings, it checks that the string length is at least that number of characters.
// For slices, arrays, and maps, validates the number of items.
func GreaterThenOrEqualValidator(p interface{}) validator.Func {
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

// lte. See https://godoc.org/github.com/go-playground/validator#hdr-Less_Than_or_Equal.
// For numbers, lte will ensure that the value is less than or equal to the parameter given.
// For strings, it checks that the string length is at most that number of characters.
// For slices, arrays, and maps, validates the number of items.
func LessThenOrEqualValidator(p interface{}) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return lenHelper(i, p, func(i, p int64) bool {
			return i <= p
		}, func(i, p uint64) bool {
			return i <= p
		}, func(i, p float64) bool {
			return xnumber.DefaultAccuracy.SmallerOrEqual(i, p)
		})
	}
}

// min, max.
// Combine GreaterThenOrEqualValidator and LessThenOrEqualValidator with And.
func LengthRangeValidator(min, max interface{}) validator.Func {
	return And(GreaterThenOrEqualValidator(min), LessThenOrEqualValidator(max))
}

// min, max.
// Combine GreaterThenOrEqualValidator and LessThenOrEqualValidator with Or.
func LengthOutOfRangeValidator(min, max interface{}) validator.Func {
	return Or(GreaterThenOrEqualValidator(max), LessThenOrEqualValidator(min))
}

// oneof
// For strings, ints, uints, and floats, oneof will ensure that the value is one of the values in the parameter.
func OneofValidator(ps ...interface{}) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		iv, err := xreflect.IufsOf(i)
		if err != nil {
			return false
		}

		for _, p := range ps {
			switch iv.Flag() {
			case xreflect.Int: // int
				p, ok := xreflect.GetInt(p)
				if ok && iv.Int() == p {
					return true
				}
			case xreflect.Uint: // uint
				p, ok := xreflect.GetUint(p)
				if ok && iv.Uint() == p {
					return true
				}
			case xreflect.String: // string
				p, ok := xreflect.GetString(p)
				if ok && iv.String() == p {
					return true
				}
			case xreflect.Float: // float
				p, ok := xreflect.GetFloat(p)
				if ok && xnumber.DefaultAccuracy.Equal(iv.Float(), p) {
					return true
				}
			}
		}
		return false
	}
}
