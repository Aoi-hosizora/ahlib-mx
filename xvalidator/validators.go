package xvalidator

import (
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/Aoi-hosizora/ahlib/xreflect"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"time"
)

// =================
// regexp & datetime
// =================

// ParamRegexpValidator represents parameterized regexp validator, just like `regexp: xxx`. For more regexps, see xvalidator.regexps package and
// https://github.com/go-playground/validator/blob/master/regexes.go.
func ParamRegexpValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		regexpParam := fl.Param() // param
		i := fl.Field().Interface()
		text, ok := i.(string)
		if !ok {
			return false
		}
		re, err := regexp.Compile(regexpParam)
		if err != nil {
			return false // return false
		}
		return re.MatchString(text)
	}
}

// RegexpValidator represents regexp validator using given regexp.
func RegexpValidator(re *regexp.Regexp) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		text, ok := i.(string)
		if !ok {
			return false
		}
		return re.MatchString(text)
	}
}

// DateTimeValidator represents datetime validator using given layout.
func DateTimeValidator(layout string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		text, ok := i.(string)
		if !ok {
			return false
		}
		_, err := time.Parse(layout, text)
		return err == nil
	}
}

// The following validators are referenced from https://github.com/go-playground/validator/blob/master/baked_in.go.

// ==============
// and & or & not
// ==============

const (
	panicNilValidatorFunc = "xvalidator: nil validator function"
)

// And represents the intersection of multiple validators, just like ',' in validator tag.
func And(fns ...validator.Func) validator.Func {
	for _, fn := range fns {
		if fn == nil {
			panic(panicNilValidatorFunc)
		}
	}
	return func(fl validator.FieldLevel) bool {
		for _, fn := range fns {
			if !fn(fl) {
				return false
			}
		}
		return true
	}
}

// Or represents the union of multiple validators, just like '|' in validator tag.
// See https://godoc.org/github.com/go-playground/validator#hdr-Or_Operator.
func Or(fns ...validator.Func) validator.Func {
	for _, fn := range fns {
		if fn == nil {
			panic(panicNilValidatorFunc)
		}
	}
	return func(fl validator.FieldLevel) bool {
		for _, fn := range fns {
			if fn(fl) {
				return true
			}
		}
		return false
	}
}

// Not represents the inverse result of given validator.
func Not(fn validator.Func) validator.Func {
	if fn == nil {
		panic(panicNilValidatorFunc)
	}
	return func(fl validator.FieldLevel) bool {
		return !fn(fl)
	}
}

// ==========
// validators
// ==========

// EqualValidator represents `eq` validator tag. See https://godoc.org/github.com/go-playground/validator#hdr-Equals.
// For strings & numbers, eq will ensure that the value is equal to the parameter given.
// For slices, arrays, and maps, validates the number of items.
func EqualValidator(p interface{}) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		eq, valid := eqHelper(i, p)
		return valid && eq
	}
}

// NotEqualValidator represents `ne` validator tag. See https://godoc.org/github.com/go-playground/validator#hdr-Not_Equal.
// For strings & numbers, ne will ensure that the value is not equal to the parameter given.
// For slices, arrays, and maps, validates the number of items.
func NotEqualValidator(p interface{}) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		eq, valid := eqHelper(i, p)
		return valid && !eq
	}
}

// LenValidator represents `len` validator tag. See https://godoc.org/github.com/go-playground/validator#hdr-Length.
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
			return xnumber.EqualInAccuracy(i, p)
		})
	}
}

// GreaterThenValidator represents `gt` validator tag. See https://godoc.org/github.com/go-playground/validator#hdr-Greater_Than.
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
			return xnumber.GreaterInAccuracy(i, p)
		})
	}
}

// LessThenValidator represents `lt` validator tag. See https://godoc.org/github.com/go-playground/validator#hdr-Less_Than.
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
			return xnumber.LessInAccuracy(i, p)
		})
	}
}

// GreaterThenOrEqualValidator represents `gte` validator tag. See https://godoc.org/github.com/go-playground/validator#hdr-Greater_Than_or_Equal.
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
			return xnumber.GreaterOrEqualInAccuracy(i, p)
		})
	}
}

// LessThenOrEqualValidator represents `lte` validator tag. See https://godoc.org/github.com/go-playground/validator#hdr-Less_Than_or_Equal.
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
			return xnumber.LessOrEqualInAccuracy(i, p)
		})
	}
}

// LengthInRangeValidator represents `min,max` validator tag, equals to combine GreaterThenOrEqualValidator and LessThenOrEqualValidator with And.
func LengthInRangeValidator(min, max interface{}) validator.Func {
	return And(GreaterThenOrEqualValidator(min), LessThenOrEqualValidator(max)) // min <= p && p <= max
}

// LengthOutOfRangeValidator represents `min|max` validator tag, equals to combine GreaterThenOrEqualValidator and LessThenOrEqualValidator with Or.
func LengthOutOfRangeValidator(min, max interface{}) validator.Func {
	return Or(GreaterThenValidator(max), LessThenValidator(min)) // p <= min || max <= p
}

// OneofValidator represents `oneof` validator tag. See https://godoc.org/github.com/go-playground/validator#hdr-One_Of.
// For strings, ints, uints, and floats, oneof will ensure that the value is one of the values in the parameter.
func OneofValidator(ps ...interface{}) validator.Func {
	return func(fl validator.FieldLevel) bool {
		i := fl.Field().Interface()
		return oneofHelper(i, ps)
	}
}

// =======
// helpers
// =======

// eqHelper is a helper function for equality used for EqualValidator and NotEqualValidator.
// For numbers & strings, it validates the value.
// For slices, arrays, and maps, it validates the length.
func eqHelper(i, p interface{}) (eq bool, valid bool) {
	iv, pv := reflect.ValueOf(i), reflect.ValueOf(p)
	ik, pk := iv.Kind(), pv.Kind()
	switch {
	case xreflect.IsIntKind(ik):
		return xreflect.IsIntKind(pk) && iv.Int() == pv.Int(), true
	case xreflect.IsUintKind(ik):
		return xreflect.IsUintKind(pk) && iv.Uint() == pv.Uint(), true
	case xreflect.IsFloatKind(ik):
		return xreflect.IsFloatKind(pk) && iv.Float() == pv.Float(), true
	case ik == reflect.Bool:
		return pk == reflect.Bool && iv.Bool() == pv.Bool(), true
	case ik == reflect.String:
		return pk == reflect.String && iv.String() == pv.String(), true
	case ik == reflect.Slice || ik == reflect.Array || ik == reflect.Map:
		return pk == reflect.Int && int64(iv.Len()) == pv.Int(), true
	}
	return false, false
}

// lenHelper is a helper function for length comparison used for LenValidator, GreaterThenValidator, LessThenValidator, GreaterThenOrEqualValidator and LessThenOrEqualValidator.
// For numbers, it validates the value.
// For strings, it validates the length of string.
// For slices, arrays, and maps, it validates the length.
func lenHelper(i, p interface{}, fi func(i, p int64) bool, fu func(i, p uint64) bool, ff func(i, p float64) bool) bool {
	iv, pv := reflect.ValueOf(i), reflect.ValueOf(p)
	ik, pk := iv.Kind(), pv.Kind()
	switch {
	case xreflect.IsIntKind(ik):
		return xreflect.IsIntKind(pk) && fi(iv.Int(), pv.Int())
	case xreflect.IsUintKind(ik):
		return xreflect.IsUintKind(pk) && fu(iv.Uint(), pv.Uint())
	case xreflect.IsFloatKind(ik):
		return xreflect.IsFloatKind(pk) && ff(iv.Float(), pv.Float())
	case ik == reflect.Bool:
		return pk == reflect.Bool && fi(int64(xnumber.Bool(iv.Bool())), int64(xnumber.Bool(pv.Bool())))
	case ik == reflect.String:
		return pk == reflect.Int && fi(int64(len([]rune(iv.String()))), pv.Int())
	case ik == reflect.Slice || ik == reflect.Array || ik == reflect.Map:
		return pk == reflect.Int && fi(int64(iv.Len()), pv.Int())
	}
	return false
}

// oneofHelper is a helper function for oneof used for OneofValidator.
// For numbers & strings, it validates the value.
func oneofHelper(i interface{}, ps []interface{}) bool {
	iv := reflect.ValueOf(i)
	ik := iv.Kind()
	switch {
	case xreflect.IsIntKind(ik):
		for _, p := range ps {
			if pv := reflect.ValueOf(p); xreflect.IsIntKind(pv.Kind()) && iv.Int() == pv.Int() {
				return true
			}
		}
		return false
	case xreflect.IsUintKind(ik):
		for _, p := range ps {
			if pv := reflect.ValueOf(p); xreflect.IsUintKind(pv.Kind()) && iv.Uint() == pv.Uint() {
				return true
			}
		}
		return false
	case xreflect.IsFloatKind(ik):
		for _, p := range ps {
			if pv := reflect.ValueOf(p); xreflect.IsFloatKind(pv.Kind()) && xnumber.EqualInAccuracy(iv.Float(), pv.Float()) {
				return true
			}
		}
		return false
	case ik == reflect.Bool:
		for _, p := range ps {
			if pv := reflect.ValueOf(p); pv.Kind() == reflect.Bool && iv.Bool() == pv.Bool() {
				return true
			}
		}
		return false
	case ik == reflect.String:
		for _, p := range ps {
			if pv := reflect.ValueOf(p); pv.Kind() == reflect.Bool && iv.String() == pv.String() {
				return true
			}
		}
		return false
	}

	// other types
	return false
}
