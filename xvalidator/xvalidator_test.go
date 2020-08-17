package xvalidator

import (
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestRequired(t *testing.T) {
	val := validator.New()
	type s struct {
		Int int `validate:"required"`
	}

	s1 := &s{}
	s2 := &s{Int: 1}
	xtesting.True(t, ValidationRequiredError(val.Struct(s1)))
	xtesting.False(t, ValidationRequiredError(val.Struct(s2)))
}

type s struct {
	Int    int     `validate:"int"`
	Uint   uint    `validate:"uint"`
	Float  float32 `validate:"float"`
	String string  `validate:"string"`
	Slice  []int   `validate:"slice"`
}

func TestEqual(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", EqualValidator(5))
	_ = val.RegisterValidation("uint", EqualValidator(uint(5)))
	_ = val.RegisterValidation("float", EqualValidator(5.0))
	_ = val.RegisterValidation("string", EqualValidator("5"))
	_ = val.RegisterValidation("slice", EqualValidator(5))

	s1 := &s{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		String: "5",
		Slice:  []int{5, 5, 5, 5, 5},
	}
	s2 := &s{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		String: "4",
		Slice:  []int{5, 5, 5, 5},
	}
	xtesting.Nil(t, show(val.Struct(s1)))
	xtesting.NotNil(t, show(val.Struct(s2)))
}

func TestNotEqual(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", NotEqualValidator(5))
	_ = val.RegisterValidation("uint", NotEqualValidator(uint(5)))
	_ = val.RegisterValidation("float", NotEqualValidator(5.0))
	_ = val.RegisterValidation("string", NotEqualValidator("5"))
	_ = val.RegisterValidation("slice", NotEqualValidator(5))

	s1 := &s{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		String: "5",
		Slice:  []int{5, 5, 5, 5, 5},
	}
	s2 := &s{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		String: "4",
		Slice:  []int{5, 5, 5, 5},
	}
	xtesting.NotNil(t, show(val.Struct(s1)))
	xtesting.Nil(t, show(val.Struct(s2)))
}

func TestLen(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", LenValidator(5))
	_ = val.RegisterValidation("uint", LenValidator(uint(5)))
	_ = val.RegisterValidation("float", LenValidator(5.0))
	_ = val.RegisterValidation("string", LenValidator(5))
	_ = val.RegisterValidation("slice", LenValidator(5))

	s1 := &s{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		String: "55555",
		Slice:  []int{5, 5, 5, 5, 5},
	}
	s2 := &s{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		String: "4444",
		Slice:  []int{5, 5, 5, 5},
	}
	xtesting.Nil(t, show(val.Struct(s1)))
	xtesting.NotNil(t, show(val.Struct(s2)))
}

func TestGreaterThen(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", GreaterThenValidator(5))
	_ = val.RegisterValidation("uint", GreaterThenValidator(uint(5)))
	_ = val.RegisterValidation("float", GreaterThenValidator(5.0))
	_ = val.RegisterValidation("string", GreaterThenValidator(5))
	_ = val.RegisterValidation("slice", GreaterThenValidator(5))

	s1 := &s{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		String: "55555",
		Slice:  []int{5, 5, 5, 5, 5},
	}
	s2 := &s{
		Int:    6,
		Uint:   6,
		Float:  6.0,
		String: "666666",
		Slice:  []int{6, 6, 6, 6, 6, 6},
	}
	xtesting.NotNil(t, show(val.Struct(s1)))
	xtesting.Nil(t, show(val.Struct(s2)))
}

func TestGreaterThenOrEqual(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", GreaterThenValidator(5))
	_ = val.RegisterValidation("uint", GreaterThenValidator(uint(5)))
	_ = val.RegisterValidation("float", GreaterThenValidator(5.0))
	_ = val.RegisterValidation("string", GreaterThenValidator(5))
	_ = val.RegisterValidation("slice", GreaterThenValidator(5))

	s1 := &s{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		String: "55555",
		Slice:  []int{5, 5, 5, 5, 5},
	}
	s2 := &s{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		String: "4444",
		Slice:  []int{4, 4, 4, 4},
	}
	xtesting.Nil(t, show(val.Struct(s1)))
	xtesting.NotNil(t, show(val.Struct(s2)))
}

func TestLessThen(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", LessThenValidator(5))
	_ = val.RegisterValidation("uint", LessThenValidator(uint(5)))
	_ = val.RegisterValidation("float", LessThenValidator(5.0))
	_ = val.RegisterValidation("string", LessThenValidator(5))
	_ = val.RegisterValidation("slice", LessThenValidator(5))

	s1 := &s{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		String: "55555",
		Slice:  []int{5, 5, 5, 5, 5},
	}
	s2 := &s{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		String: "4444",
		Slice:  []int{4, 4, 4, 4},
	}
	xtesting.NotNil(t, show(val.Struct(s1)))
	xtesting.Nil(t, show(val.Struct(s2)))
}

func TestLessThenOrEqual(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", LessThenOrEqualValidator(5))
	_ = val.RegisterValidation("uint", LessThenOrEqualValidator(uint(5)))
	_ = val.RegisterValidation("float", LessThenOrEqualValidator(5.0))
	_ = val.RegisterValidation("string", LessThenOrEqualValidator(5))
	_ = val.RegisterValidation("slice", LessThenOrEqualValidator(5))

	s1 := &s{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		String: "55555",
		Slice:  []int{5, 5, 5, 5, 5},
	}
	s2 := &s{
		Int:    6,
		Uint:   6,
		Float:  6.0,
		String: "666666",
		Slice:  []int{6, 6, 6, 6, 6, 6},
	}
	xtesting.Nil(t, show(val.Struct(s1)))
	xtesting.NotNil(t, show(val.Struct(s2)))
}

func TestLengthRange(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", LengthRangeValidator(5, 6))
	_ = val.RegisterValidation("uint", LengthRangeValidator(uint(5), uint(6)))
	_ = val.RegisterValidation("float", LengthRangeValidator(5.0, 6.0))
	_ = val.RegisterValidation("string", LengthRangeValidator(5, 6))
	_ = val.RegisterValidation("slice", LengthRangeValidator(5, 6))

	s1 := &s{
		Int:    5,
		Uint:   5,
		Float:  6.0,
		String: "55555",
		Slice:  []int{6, 6, 6, 6, 6, 6},
	}
	s2 := &s{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		String: "4444",
		Slice:  []int{4, 4, 4, 4},
	}
	xtesting.Nil(t, show(val.Struct(s1)))
	xtesting.NotNil(t, show(val.Struct(s2)))
}

func TestLengthOutOfRange(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", LengthOutOfRangeValidator(5, 7))
	_ = val.RegisterValidation("uint", LengthOutOfRangeValidator(uint(5), uint(7)))
	_ = val.RegisterValidation("float", LengthOutOfRangeValidator(5.0, 7.0))
	_ = val.RegisterValidation("string", LengthOutOfRangeValidator(5, 7))
	_ = val.RegisterValidation("slice", LengthOutOfRangeValidator(5, 7))

	s1 := &s{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		String: "4444",
		Slice:  []int{4, 4, 4, 4},
	}
	s2 := &s{
		Int:    6,
		Uint:   6,
		Float:  6.0,
		String: "666666",
		Slice:  []int{6, 6, 6, 6, 6, 6},
	}
	xtesting.Nil(t, show(val.Struct(s1)))
	xtesting.NotNil(t, show(val.Struct(s2)))
}

func TestOneof(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", OneofValidator(1, 2, 3))
	_ = val.RegisterValidation("uint", OneofValidator(uint(1), uint(2), uint(3)))
	_ = val.RegisterValidation("float", OneofValidator(1.0, 2.0, 3.0))
	_ = val.RegisterValidation("string", OneofValidator("1", "2", "3"))
	_ = val.RegisterValidation("slice", LenValidator(0))

	s1 := &s{
		Int:    1,
		Uint:   2,
		Float:  3.0,
		String: "2",
		Slice:  []int{},
	}
	s2 := &s{
		Int:    4,
		Uint:   5,
		Float:  6.0,
		String: "22",
		Slice:  []int{},
	}
	xtesting.Nil(t, show(val.Struct(s1)))
	xtesting.NotNil(t, show(val.Struct(s2)))
}

func TestOtherValidator(t *testing.T) {
	type s struct {
		Regexp1 string `validate:"regexp=^[0123]+$"`
		Regexp2 string `validate:"rr"`
		Date    string `validate:"datetime"`
	}

	val := validator.New()
	_ = val.RegisterValidation("regexp", DefaultRegexpValidator())
	_ = val.RegisterValidation("rr", RegexpValidator(regexp.MustCompile(`^[abc]+$`)))
	_ = val.RegisterValidation("datetime", DateTimeValidator(time.RFC3339))

	s1 := &s{
		Regexp1: "01230123",
		Regexp2: "abcabcabc",
		Date:    time.Now().Format(time.RFC3339),
	}
	s2 := &s{
		Regexp1: "012340123",
		Regexp2: "abcadbcabc",
		Date:    time.Now().Format(time.RFC1123Z),
	}
	xtesting.Nil(t, show(val.Struct(s1)))
	xtesting.NotNil(t, show(val.Struct(s2)))
}

func show(err interface{}) interface{} {
	if err == nil {
		log.Println(nil)
		return nil
	}
	errs := err.(validator.ValidationErrors)
	if len(errs) == 0 {
		return nil
	}
	sp := strings.Builder{}
	for _, err := range errs {
		sp.WriteString(err.Field() + ":" + err.Tag() + ", ")
	}
	l := sp.String()
	log.Println(l[:len(l)-2])
	return err
}
