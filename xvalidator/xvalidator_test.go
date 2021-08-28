package xvalidator

import (
	"errors"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/Aoi-hosizora/ahlib/xtime"
	"github.com/go-playground/validator/v10"
	"regexp"
	"testing"
)

func TestIsXXXError(t *testing.T) {
	val := validator.New()
	type testStruct struct {
		Int int `validate:"required,lt=2"`
	}

	for _, tc := range []struct {
		giveErr   error
		wantReqOk bool
		wantValOk bool
	}{
		{nil, false, false},
		{errors.New("test"), false, false},
		{validator.ValidationErrors{}, false, true},
		{val.Struct(&testStruct{}), true, true},
		{val.Struct(&testStruct{Int: 0}), true, true},
		{val.Struct(&testStruct{Int: 1}), false, false},
		{val.Struct(&testStruct{Int: 3}), false, true},
	} {
		xtesting.Equal(t, IsValidationError(tc.giveErr), tc.wantValOk)
		xtesting.Equal(t, IsRequiredError(tc.giveErr), tc.wantReqOk)
	}
}

func TestRegexpAndDateTimeValidator(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("regexp", ParamRegexpValidator())
	_ = val.RegisterValidation("abc", RegexpValidator(regexp.MustCompile(`^[abc].+$`)))
	_ = val.RegisterValidation("date", DateTimeValidator(xtime.RFC3339Date))
	_ = val.RegisterValidation("datetime", DateTimeValidator(xtime.RFC3339DateTime))

	for _, tc := range []struct {
		give    interface{}
		wantErr bool
	}{
		{&struct{}{}, false},

		{&struct {
			Int int `validate:"regexp=^[abc]+$"`
		}{}, true},
		{&struct {
			String string `validate:"regexp=^[$"`
		}{}, true},
		{&struct {
			String string `validate:"regexp=^[abc]+$"`
		}{"abcd"}, true},
		{&struct {
			String string `validate:"regexp=^[abc]+$"`
		}{"abc"}, false},

		{&struct {
			Int int `validate:"abc"`
		}{}, true},
		{&struct {
			String string `validate:"abc"`
		}{"dcba"}, true},
		{&struct {
			String string `validate:"abc"`
		}{"abc"}, false},
		{&struct {
			String string `validate:"abc=dummy"`
		}{"abc"}, false},

		{&struct {
			Date     int `validate:"date"`
			DateTime int `validate:"datetime"`
		}{}, true},
		{&struct {
			Date     string `validate:"date"`
			DateTime string `validate:"datetime"`
		}{"2021/01/24", "2021-01-24T02:55:29"}, true},
		{&struct {
			Date     string `validate:"date"`
			DateTime string `validate:"datetime"`
		}{"2021-01-24", "2021-01-24T02:55:29+08:00"}, false},
		{&struct {
			Date string `validate:"date=dummy"`
		}{"2021-01-24"}, false},
	} {
		err := val.Struct(tc.give)
		if tc.wantErr {
			xtesting.NotNil(t, err)
		} else {
			xtesting.Nil(t, err)
		}
	}
}

func TestAndOr(t *testing.T) {
	val := validator.New()

	xtesting.NotPanic(t, func() { And() })
	xtesting.Panic(t, func() { And(nil, nil, nil) })
	xtesting.Panic(t, func() { And(ParamRegexpValidator(), nil, nil) })
	xtesting.NotPanic(t, func() { Or() })
	xtesting.Panic(t, func() { Or(nil, nil, nil) })
	xtesting.Panic(t, func() { Or(ParamRegexpValidator(), nil, nil) })

	_ = val.RegisterValidation("re", And(RegexpValidator(regexp.MustCompile(`^[abc].+$`)), RegexpValidator(regexp.MustCompile(`^[abc][def].+$`))))
	_ = val.RegisterValidation("time", Or(DateTimeValidator(xtime.RFC3339Date), DateTimeValidator(xtime.RFC3339DateTime)))

	type testStruct struct {
		Re   string `validate:"re"`
		Time string `validate:"time"`
	}

	for _, tc := range []struct {
		give    *testStruct
		wantErr bool
	}{
		{&testStruct{"", ""}, true},
		{&testStruct{"aaa", "2021/01/24"}, true},
		{&testStruct{"ada", "2021-01-24"}, false},
		{&testStruct{"aef", "2021-01-24T15:51:22+08:00"}, false},
	} {
		if tc.wantErr {
			xtesting.NotNil(t, val.Struct(tc.give))
		} else {
			xtesting.Nil(t, val.Struct(tc.give))
		}
	}
}

type testStruct struct {
	Int    int     `validate:"int"`
	Uint   uint    `validate:"uint"`
	Float  float32 `validate:"float"`
	Bool   bool    `validate:"bool"`
	String string  `validate:"string"`
	Slice  []int   `validate:"slice"`
}

type malformedStruct struct {
	Complex complex128 `validate:"complex"`
	Fn      func()     `validate:"fn"`
}

func getError(err error) []string {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	fields := make([]string, len(errs))
	for idx, e := range errs {
		fields[idx] = e.Tag()
	}
	return fields
}

func TestEqual(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", EqualValidator(5))
	_ = val.RegisterValidation("uint", EqualValidator(uint(5)))
	_ = val.RegisterValidation("float", EqualValidator(5.0))
	_ = val.RegisterValidation("bool", EqualValidator(true))
	_ = val.RegisterValidation("string", EqualValidator("5"))
	_ = val.RegisterValidation("slice", EqualValidator(5))
	_ = val.RegisterValidation("complex", EqualValidator(5i))
	_ = val.RegisterValidation("fn", EqualValidator(func() {}))

	s1 := &testStruct{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		Bool:   true,
		String: "5",
		Slice:  []int{5, 5, 5, 5, 5}, // 5
	}
	s2 := &testStruct{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		Bool:   false,
		String: "4",
		Slice:  []int{4, 4, 4, 4}, // 4
	}
	err1 := val.Struct(s1)
	err2 := val.Struct(s2)
	xtesting.Nil(t, err1)
	xtesting.NotNil(t, err2)
	xtesting.ElementMatch(t, getError(err2), []string{"int", "uint", "float", "bool", "string", "slice"})

	s3 := &malformedStruct{Complex: 5i, Fn: func() {}}
	s4 := &malformedStruct{Complex: 4i, Fn: func() {}}
	err3 := val.Struct(s3)
	err4 := val.Struct(s4)
	xtesting.NotNil(t, err3)
	xtesting.NotNil(t, err4)
	xtesting.ElementMatch(t, getError(err3), []string{"complex", "fn"})
	xtesting.ElementMatch(t, getError(err4), []string{"complex", "fn"})
}

func TestNotEqual(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", NotEqualValidator(5))
	_ = val.RegisterValidation("uint", NotEqualValidator(uint(5)))
	_ = val.RegisterValidation("float", NotEqualValidator(5.0))
	_ = val.RegisterValidation("bool", NotEqualValidator(true))
	_ = val.RegisterValidation("string", NotEqualValidator("5"))
	_ = val.RegisterValidation("slice", NotEqualValidator(5))
	_ = val.RegisterValidation("complex", NotEqualValidator(5i))
	_ = val.RegisterValidation("fn", NotEqualValidator(func() {}))

	s1 := &testStruct{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		Bool:   true,
		String: "5",
		Slice:  []int{5, 5, 5, 5, 5}, // 5
	}
	s2 := &testStruct{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		Bool:   false,
		String: "4",
		Slice:  []int{4, 4, 4, 4}, // 4
	}
	err1 := val.Struct(s1)
	err2 := val.Struct(s2)
	xtesting.NotNil(t, err1)
	xtesting.Nil(t, err2)
	xtesting.ElementMatch(t, getError(err1), []string{"int", "uint", "float", "bool", "string", "slice"})

	s3 := &malformedStruct{Complex: 5i, Fn: func() {}}
	s4 := &malformedStruct{Complex: 4i, Fn: func() {}}
	err3 := val.Struct(s3)
	err4 := val.Struct(s4)
	xtesting.NotNil(t, err3)
	xtesting.NotNil(t, err4)
	xtesting.ElementMatch(t, getError(err3), []string{"complex", "fn"})
	xtesting.ElementMatch(t, getError(err4), []string{"complex", "fn"})
}

func TestLen(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", LenValidator(5))
	_ = val.RegisterValidation("uint", LenValidator(uint(5)))
	_ = val.RegisterValidation("float", LenValidator(5.0))
	_ = val.RegisterValidation("bool", LenValidator(true))
	_ = val.RegisterValidation("string", LenValidator(5))
	_ = val.RegisterValidation("slice", LenValidator(5))
	_ = val.RegisterValidation("complex", LenValidator(5i))
	_ = val.RegisterValidation("fn", LenValidator(func() {}))

	s1 := &testStruct{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		Bool:   true,
		String: "55555",              // 5
		Slice:  []int{5, 5, 5, 5, 5}, // 5
	}
	s2 := &testStruct{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		Bool:   false,
		String: "4444",            // 4
		Slice:  []int{4, 4, 4, 4}, // 4
	}
	err1 := val.Struct(s1)
	err2 := val.Struct(s2)
	xtesting.Nil(t, err1)
	xtesting.NotNil(t, err2)
	xtesting.ElementMatch(t, getError(err2), []string{"int", "uint", "float", "bool", "string", "slice"})

	s3 := &malformedStruct{Complex: 5i, Fn: func() {}}
	s4 := &malformedStruct{Complex: 4i, Fn: func() {}}
	err3 := val.Struct(s3)
	err4 := val.Struct(s4)
	xtesting.NotNil(t, err3)
	xtesting.NotNil(t, err4)
	xtesting.ElementMatch(t, getError(err3), []string{"complex", "fn"})
	xtesting.ElementMatch(t, getError(err4), []string{"complex", "fn"})
}

func TestGreaterThen(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", GreaterThenValidator(5))
	_ = val.RegisterValidation("uint", GreaterThenValidator(uint(5)))
	_ = val.RegisterValidation("float", GreaterThenValidator(5.0))
	_ = val.RegisterValidation("bool", GreaterThenValidator(false))
	_ = val.RegisterValidation("string", GreaterThenValidator(5))
	_ = val.RegisterValidation("slice", GreaterThenValidator(5))
	_ = val.RegisterValidation("complex", GreaterThenValidator(5i))
	_ = val.RegisterValidation("fn", GreaterThenValidator(func() {}))

	s1 := &testStruct{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		Bool:   false,
		String: "55555",              // 5
		Slice:  []int{5, 5, 5, 5, 5}, // 5
	}
	s2 := &testStruct{
		Int:    6,
		Uint:   6,
		Float:  6.0,
		Bool:   true,
		String: "666666",                // 6
		Slice:  []int{6, 6, 6, 6, 6, 6}, // 6
	}
	err1 := val.Struct(s1)
	err2 := val.Struct(s2)
	xtesting.NotNil(t, err1)
	xtesting.Nil(t, err2)
	xtesting.ElementMatch(t, getError(err1), []string{"int", "uint", "float", "bool", "string", "slice"})

	s3 := &malformedStruct{Complex: 5i, Fn: func() {}}
	s4 := &malformedStruct{Complex: 6i, Fn: func() {}}
	err3 := val.Struct(s3)
	err4 := val.Struct(s4)
	xtesting.NotNil(t, err3)
	xtesting.NotNil(t, err4)
	xtesting.ElementMatch(t, getError(err3), []string{"complex", "fn"})
	xtesting.ElementMatch(t, getError(err4), []string{"complex", "fn"})
}

func TestGreaterThenOrEqual(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", GreaterThenOrEqualValidator(5))
	_ = val.RegisterValidation("uint", GreaterThenOrEqualValidator(uint(5)))
	_ = val.RegisterValidation("float", GreaterThenOrEqualValidator(5.0))
	_ = val.RegisterValidation("bool", GreaterThenOrEqualValidator(false))
	_ = val.RegisterValidation("string", GreaterThenOrEqualValidator(5))
	_ = val.RegisterValidation("slice", GreaterThenOrEqualValidator(5))
	_ = val.RegisterValidation("complex", GreaterThenOrEqualValidator(5i))
	_ = val.RegisterValidation("fn", GreaterThenOrEqualValidator(func() {}))

	s1 := &testStruct{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		Bool:   false,
		String: "55555",              // 5
		Slice:  []int{5, 5, 5, 5, 5}, // 5
	}
	s2 := &testStruct{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		Bool:   true,
		String: "4444",            // 4
		Slice:  []int{4, 4, 4, 4}, // 4
	}
	err1 := val.Struct(s1)
	err2 := val.Struct(s2)
	xtesting.Nil(t, err1)
	xtesting.NotNil(t, err2)
	xtesting.ElementMatch(t, getError(err2), []string{"int", "uint", "float", "string", "slice"})

	s3 := &malformedStruct{Complex: 5i, Fn: func() {}}
	s4 := &malformedStruct{Complex: 4i, Fn: func() {}}
	err3 := val.Struct(s3)
	err4 := val.Struct(s4)
	xtesting.NotNil(t, err3)
	xtesting.NotNil(t, err4)
	xtesting.ElementMatch(t, getError(err3), []string{"complex", "fn"})
	xtesting.ElementMatch(t, getError(err4), []string{"complex", "fn"})
}

func TestLessThen(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", LessThenValidator(5))
	_ = val.RegisterValidation("uint", LessThenValidator(uint(5)))
	_ = val.RegisterValidation("float", LessThenValidator(5.0))
	_ = val.RegisterValidation("bool", LessThenValidator(true))
	_ = val.RegisterValidation("string", LessThenValidator(5))
	_ = val.RegisterValidation("slice", LessThenValidator(5))
	_ = val.RegisterValidation("complex", LessThenValidator(5i))
	_ = val.RegisterValidation("fn", LessThenValidator(func() {}))

	s1 := &testStruct{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		Bool:   true,
		String: "55555",              // 5
		Slice:  []int{5, 5, 5, 5, 5}, // 5
	}
	s2 := &testStruct{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		Bool:   false,
		String: "4444",            // 4
		Slice:  []int{4, 4, 4, 4}, // 4
	}
	err1 := val.Struct(s1)
	err2 := val.Struct(s2)
	xtesting.NotNil(t, err1)
	xtesting.Nil(t, err2)
	xtesting.ElementMatch(t, getError(err1), []string{"int", "uint", "float", "bool", "string", "slice"})

	s3 := &malformedStruct{Complex: 5i, Fn: func() {}}
	s4 := &malformedStruct{Complex: 4i, Fn: func() {}}
	err3 := val.Struct(s3)
	err4 := val.Struct(s4)
	xtesting.NotNil(t, err3)
	xtesting.NotNil(t, err4)
	xtesting.ElementMatch(t, getError(err3), []string{"complex", "fn"})
	xtesting.ElementMatch(t, getError(err4), []string{"complex", "fn"})
}

func TestLessThenOrEqual(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", LessThenOrEqualValidator(5))
	_ = val.RegisterValidation("uint", LessThenOrEqualValidator(uint(5)))
	_ = val.RegisterValidation("float", LessThenOrEqualValidator(5.0))
	_ = val.RegisterValidation("bool", LessThenOrEqualValidator(true))
	_ = val.RegisterValidation("string", LessThenOrEqualValidator(5))
	_ = val.RegisterValidation("slice", LessThenOrEqualValidator(5))
	_ = val.RegisterValidation("complex", LessThenOrEqualValidator(5i))
	_ = val.RegisterValidation("fn", LessThenOrEqualValidator(func() {}))

	s1 := &testStruct{
		Int:    5,
		Uint:   5,
		Float:  5.0,
		Bool:   true,
		String: "55555",              // 5
		Slice:  []int{5, 5, 5, 5, 5}, // 5
	}
	s2 := &testStruct{
		Int:    6,
		Uint:   6,
		Float:  6.0,
		Bool:   false,
		String: "666666",                // 6
		Slice:  []int{6, 6, 6, 6, 6, 6}, // 6
	}
	err1 := val.Struct(s1)
	err2 := val.Struct(s2)
	xtesting.Nil(t, err1)
	xtesting.NotNil(t, err2)
	xtesting.ElementMatch(t, getError(err2), []string{"int", "uint", "float", "string", "slice"})

	s3 := &malformedStruct{Complex: 5i, Fn: func() {}}
	s4 := &malformedStruct{Complex: 6i, Fn: func() {}}
	err3 := val.Struct(s3)
	err4 := val.Struct(s4)
	xtesting.NotNil(t, err3)
	xtesting.NotNil(t, err4)
	xtesting.ElementMatch(t, getError(err3), []string{"complex", "fn"})
	xtesting.ElementMatch(t, getError(err4), []string{"complex", "fn"})
}

func TestLengthRange(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", LengthInRangeValidator(5, 6))
	_ = val.RegisterValidation("uint", LengthInRangeValidator(uint(5), uint(6)))
	_ = val.RegisterValidation("float", LengthInRangeValidator(5.0, 6.0))
	_ = val.RegisterValidation("bool", LengthInRangeValidator(true, true))
	_ = val.RegisterValidation("string", LengthInRangeValidator(5, 6))
	_ = val.RegisterValidation("slice", LengthInRangeValidator(5, 6))
	_ = val.RegisterValidation("complex", LengthInRangeValidator(5i, 6i))
	_ = val.RegisterValidation("fn", LengthInRangeValidator(func() {}, func() {}))

	s1 := &testStruct{
		Int:    6,
		Uint:   6,
		Float:  6.0,
		Bool:   true,
		String: "666666",                // 6
		Slice:  []int{6, 6, 6, 6, 6, 6}, // 6
	}
	s2 := &testStruct{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		Bool:   false,
		String: "4444",            // 4
		Slice:  []int{4, 4, 4, 4}, // 4
	}
	err1 := val.Struct(s1)
	err2 := val.Struct(s2)
	xtesting.Nil(t, err1)
	xtesting.NotNil(t, err2)
	xtesting.ElementMatch(t, getError(err2), []string{"int", "uint", "float", "bool", "string", "slice"})

	s3 := &malformedStruct{Complex: 6i, Fn: func() {}}
	s4 := &malformedStruct{Complex: 4i, Fn: func() {}}
	err3 := val.Struct(s3)
	err4 := val.Struct(s4)
	xtesting.NotNil(t, err3)
	xtesting.NotNil(t, err4)
	xtesting.ElementMatch(t, getError(err3), []string{"complex", "fn"})
	xtesting.ElementMatch(t, getError(err4), []string{"complex", "fn"})
}

func TestLengthOutOfRange(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", LengthOutOfRangeValidator(5, 7))
	_ = val.RegisterValidation("uint", LengthOutOfRangeValidator(uint(5), uint(7)))
	_ = val.RegisterValidation("float", LengthOutOfRangeValidator(5.0, 7.0))
	_ = val.RegisterValidation("bool", LengthInRangeValidator(true, true))
	_ = val.RegisterValidation("string", LengthOutOfRangeValidator(5, 7))
	_ = val.RegisterValidation("slice", LengthOutOfRangeValidator(5, 7))
	_ = val.RegisterValidation("complex", LengthOutOfRangeValidator(5i, 6i))
	_ = val.RegisterValidation("fn", LengthOutOfRangeValidator(func() {}, func() {}))

	s1 := &testStruct{
		Int:    4,
		Uint:   4,
		Float:  4.0,
		Bool:   true,
		String: "4444",            // 4
		Slice:  []int{4, 4, 4, 4}, // 4
	}
	s2 := &testStruct{
		Int:    6,
		Uint:   6,
		Float:  6.0,
		Bool:   false,
		String: "666666",                // 6
		Slice:  []int{6, 6, 6, 6, 6, 6}, // 6
	}
	err1 := val.Struct(s1)
	err2 := val.Struct(s2)
	xtesting.Nil(t, err1)
	xtesting.NotNil(t, err2)
	xtesting.ElementMatch(t, getError(err2), []string{"int", "uint", "float", "bool", "string", "slice"})

	s3 := &malformedStruct{Complex: 4i, Fn: func() {}}
	s4 := &malformedStruct{Complex: 6i, Fn: func() {}}
	err3 := val.Struct(s3)
	err4 := val.Struct(s4)
	xtesting.NotNil(t, err3)
	xtesting.NotNil(t, err4)
	xtesting.ElementMatch(t, getError(err3), []string{"complex", "fn"})
	xtesting.ElementMatch(t, getError(err4), []string{"complex", "fn"})
}

func TestOneof(t *testing.T) {
	val := validator.New()
	_ = val.RegisterValidation("int", OneofValidator(1, 2, 3))
	_ = val.RegisterValidation("uint", OneofValidator(uint(1), uint(2), uint(3)))
	_ = val.RegisterValidation("float", OneofValidator(1.0, 2.0, 3.0))
	_ = val.RegisterValidation("bool", OneofValidator(true))
	_ = val.RegisterValidation("string", OneofValidator("1", "2", "3"))
	_ = val.RegisterValidation("slice", LenValidator(0))
	_ = val.RegisterValidation("complex", OneofValidator(1i, 2i, 3i))
	_ = val.RegisterValidation("fn", OneofValidator(func() {}))

	s1 := &testStruct{
		Int:    1,
		Uint:   2,
		Float:  3.0,
		Bool:   true,
		String: "2",
		Slice:  []int{},
	}
	s2 := &testStruct{
		Int:    4,
		Uint:   5,
		Float:  6.0,
		Bool:   false,
		String: "4",
		Slice:  []int{},
	}
	err1 := val.Struct(s1)
	err2 := val.Struct(s2)
	xtesting.Nil(t, err1)
	xtesting.NotNil(t, err2)
	xtesting.ElementMatch(t, getError(err2), []string{"int", "uint", "float", "bool", "string"})

	s3 := &malformedStruct{Complex: 1i, Fn: func() {}}
	s4 := &malformedStruct{Complex: 4i, Fn: func() {}}
	err3 := val.Struct(s3)
	err4 := val.Struct(s4)
	xtesting.NotNil(t, err3)
	xtesting.NotNil(t, err4)
	xtesting.ElementMatch(t, getError(err3), []string{"complex", "fn"})
	xtesting.ElementMatch(t, getError(err4), []string{"complex", "fn"})
}
