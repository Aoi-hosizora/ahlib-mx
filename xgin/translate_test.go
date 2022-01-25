package xgin

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Aoi-hosizora/ahlib-web/xvalidator"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

type translatableError string

func (t translatableError) Error() string {
	return string(t)
}

func (t translatableError) Translate() (map[string]string, bool) {
	return map[string]string{"_": string(t)}, true
}

func TestTranslateBindingError(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	val, _ := GetValidatorEngine()
	xvalidator.UseTagAsFieldName(val, "json")
	trans, _ := GetValidatorTranslator(xvalidator.EnLocaleTranslator(), xvalidator.EnTranslationRegisterFunc())
	newVal := xvalidator.NewCustomStructValidator()
	newVal.SetMessageTagName("message")
	newVal.SetValidatorTagName("binding")
	// xvalidator.UseTagAsFieldName(newVal.ValidateEngine(), "json")
	newVal.SetFieldNameTag("json")

	app := gin.New()
	type testStruct struct {
		Str string `json:"str" form:"str" binding:"required" message:"required|str should be not null and not empty"`
		Int int32  `json:"int" form:"int" binding:"required" message:"required|int should be not null and not zero"`
	}
	respond := func(c *gin.Context, code int, details map[string]string) {
		if code == 200 {
			c.JSON(200, gin.H{"success": true})
		} else {
			c.JSON(code, gin.H{"success": false, "details": details})
		}
	}

	app.POST("/body", func(c *gin.Context) {
		opts := make([]TranslateOption, 0)
		var ptr interface{} = &testStruct{}
		if c.Query("useInvalidType") == "true" {
			ptr = 0
		}
		if c.Query("useTrans") == "true" {
			opts = append(opts, WithUtTranslator(trans))
		}
		if c.Query("useCustom") == "true" {
			motoVal := binding.Validator
			binding.Validator = newVal
			defer func() { binding.Validator = motoVal }()
		}
		if err := c.ShouldBind(ptr); err != nil {
			if result, need4xx := TranslateBindingError(err, opts...); need4xx {
				respond(c, 400, result)
			} else {
				respond(c, 500, result)
			}
		} else {
			respond(c, 200, nil)
		}
	})
	app.POST("/id/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			if c.Query("useNumError") != "true" {
				err = NewRouterDecodeError("id", idStr, err, "")
			}
		} else if id <= 0 {
			err = NewRouterDecodeError("id", idStr, err, "should be larger then zero")
		}
		if rErr, ok := err.(*RouterDecodeError); ok {
			_ = rErr.Error()
			if c.Query("ignoreField") == "true" {
				rErr.Field = ""
			}
			if c.Query("ignoreMessage") == "true" {
				rErr.Message = ""
			}
		}
		if err != nil {
			if result, need4xx := TranslateBindingError(err); need4xx {
				respond(c, 400, result)
			} else {
				respond(c, 500, result)
			}
		} else {
			respond(c, 200, nil)
		}
	})

	server := &http.Server{Addr: ":12345", Handler: app}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

	// normal
	for _, tc := range []struct {
		giveRoute string
		giveBody  string
		giveQuery string
		wantCode  int
		wantMap   map[string]interface{}
	}{
		// io eof
		{"body", ``, "", 400,
			map[string]interface{}{"__decode": "requested json has an invalid syntax at position -1"},
		},
		{"body", `{`, "", 400,
			map[string]interface{}{"__decode": "requested json has an invalid syntax at position -1"},
		},
		// json invalid unmarshal
		{"body", `{}`, "useInvalidType=true", 500, nil},
		// json syntax
		{"body", `{"str": a, "int": 0}`, "useTrans=true", 400,
			map[string]interface{}{"__decode": "requested json has an invalid syntax at position 9"},
		},
		{"body", `{"str": "a, "int": 0}`, "useTrans=true", 400,
			map[string]interface{}{"__decode": "requested json has an invalid syntax at position 14"},
		},
		// json type
		{"body", `{"str": 0, "int": 0}`, "", 400,
			map[string]interface{}{"__decode": "type of 'number' in 'str' mismatches with required 'string'"},
		},
		{"body", `{"str": "", "int": ""}`, "", 400,
			map[string]interface{}{"__decode": "type of 'string' in 'int' mismatches with required 'int32'"},
		},
		{"body", `{"str": "abc", "int": 999999999999999999999999999999}`, "", 400,
			map[string]interface{}{"__decode": "type of 'number 999999999999999999999999999999' in 'int' mismatches with required 'int32'"},
		},
		{"body", `{"str": "abc", "int": 3.14}`, "", 400,
			map[string]interface{}{"__decode": "type of 'number 3.14' in 'int' mismatches with required 'int32'"},
		},
		// validator
		{"body", `{}`, "", 400,
			map[string]interface{}{"str": "Field validation for 'str' failed on the 'required' tag", "int": "Field validation for 'int' failed on the 'required' tag"},
		},
		{"body", `{}`, "useTrans=true", 400,
			map[string]interface{}{"str": "str is a required field", "int": "int is a required field"},
		},
		{"body", `{"str": "", "int": 0}`, "", 400,
			map[string]interface{}{"str": "Field validation for 'str' failed on the 'required' tag", "int": "Field validation for 'int' failed on the 'required' tag"},
		},
		{"body", `{"str": "", "int": 0}`, "useTrans=true", 400,
			map[string]interface{}{"str": "str is a required field", "int": "int is a required field"},
		},
		// xvalidator required
		{"body", `{}`, "useCustom=true", 400,
			map[string]interface{}{"str": "str should be not null and not empty", "int": "int should be not null and not zero"},
		},
		{"body", `{}`, "useCustom=true&useTrans=true", 400,
			map[string]interface{}{"str": "str should be not null and not empty", "int": "int should be not null and not zero"},
		},
		{"body", `{"str": "", "int": 0}`, "useCustom=true", 400,
			map[string]interface{}{"str": "str should be not null and not empty", "int": "int should be not null and not zero"},
		},
		{"body", `{"str": "", "int": 0}`, "useCustom=true&useTrans=true", 400,
			map[string]interface{}{"str": "str should be not null and not empty", "int": "int should be not null and not zero"},
		},
		// ok
		{"body", `{"str": "abc", "int": 1}`, "", 200, nil},
		// strconv number error
		{"id/a", ``, "useNumError=true", 400,
			map[string]interface{}{"router parameter": "router parameter must be a number"},
		},
		{"id/3.14", ``, "useNumError=true", 400,
			map[string]interface{}{"router parameter": "router parameter must be a number"},
		},
		{"id/999999999999999999999999999999", `useNumError=true`, "useNumError=true", 400,
			map[string]interface{}{"router parameter": "router parameter is out of range"},
		},
		// router decode error
		{"id/a", ``, "", 400,
			map[string]interface{}{"id": "router parameter id must be a number"},
		},
		{"id/3.14", ``, "", 400,
			map[string]interface{}{"id": "router parameter id must be a number"},
		},
		{"id/3.14", ``, "ignoreField=true", 400,
			map[string]interface{}{"router parameter": "router parameter must be a number"},
		},
		{"id/999999999999999999999999999999", ``, "", 400,
			map[string]interface{}{"id": "router parameter id is out of range"},
		},
		{"id/999999999999999999999999999999", ``, "ignoreMessage=true", 400,
			map[string]interface{}{"id": "router parameter id is out of range"},
		},
		{"id/0", ``, "", 400,
			map[string]interface{}{"id": "router parameter id should be larger then zero"},
		},
		{"id/0", ``, "ignoreField=true", 400,
			map[string]interface{}{"router parameter": "router parameter should be larger then zero"},
		},
		{"id/0", ``, "ignoreMessage=true", 500, nil},
		// ok
		{"id/1", ``, "", 200, nil},
	} {
		t.Run(tc.giveRoute+"_"+tc.giveBody, func(t *testing.T) {
			u := "http://localhost:12345/" + tc.giveRoute + "?" + tc.giveQuery
			req, _ := http.NewRequest("POST", u, strings.NewReader(tc.giveBody))
			req.Header.Set("Content-type", "application/json")
			client := &http.Client{}
			resp, err := client.Do(req)
			xtesting.Nil(t, err)
			bs, _ := ioutil.ReadAll(resp.Body)

			r := make(map[string]interface{})
			err = json.Unmarshal(bs, &r)
			xtesting.Nil(t, err)
			xtesting.Equal(t, resp.StatusCode, tc.wantCode)
			if resp.StatusCode != 200 && r["details"] != nil {
				xtesting.Equal(t, r["details"].(map[string]interface{}), tc.wantMap)
			}
		})
	}

	// other
	for _, tc := range []struct {
		name     string
		giveErr  error
		giveOpts []TranslateOption
		wantMap  map[string]string
		want4xx  bool
	}{
		{"nil", nil, nil, nil, false},
		{"NumError", &strconv.NumError{}, nil, nil, false},
		{"InvalidValidationError", &validator.InvalidValidationError{}, nil, nil, false},
		{"ExtraErrors", errors.New("TODO"), nil, nil, false},
		{"InvalidUnmarshalError", &json.InvalidUnmarshalError{}, []TranslateOption{WithJsonInvalidUnmarshalError(
			func(*json.InvalidUnmarshalError) (result map[string]string, need4xx bool) { return nil, true })}, nil, true},
		{"UnmarshalTypeError", &json.UnmarshalTypeError{}, []TranslateOption{WithJsonUnmarshalTypeError(
			func(*json.UnmarshalTypeError) (result map[string]string, need4xx bool) { return nil, false })}, nil, false},
		{"SyntaxError", &json.SyntaxError{}, []TranslateOption{WithJsonSyntaxError(
			func(*json.SyntaxError) (result map[string]string, need4xx bool) { return nil, false })}, nil, false},
		{"ErrUnexpectedEOF", io.ErrUnexpectedEOF, []TranslateOption{WithIoEOFError(
			func(error) (result map[string]string, need4xx bool) { return nil, false })}, nil, false},
		{"NumError", &strconv.NumError{}, []TranslateOption{WithStrconvNumErrorError(
			func(*strconv.NumError) (result map[string]string, need4xx bool) { return nil, false })}, nil, false},
		{"RouterDecodeError", &RouterDecodeError{}, []TranslateOption{WithXginRouterDecodeError(
			func(*RouterDecodeError) (result map[string]string, need4xx bool) { return nil, false })}, nil, false},
		{"InvalidValidationError", &validator.InvalidValidationError{}, []TranslateOption{WithValidatorInvalidTypeError(
			func(*validator.InvalidValidationError) (result map[string]string, need4xx bool) { return nil, true })}, nil, true},
		{"ValidationErrors", validator.ValidationErrors{}, []TranslateOption{WithValidatorFieldsError(
			func(validator.ValidationErrors, xvalidator.UtTranslator) (result map[string]string, need4xx bool) {
				return nil, false
			})}, nil, false},
		{"ValidateFieldsError", &xvalidator.ValidateFieldsError{}, []TranslateOption{WithXvalidatorValidateFieldsError(
			func(*xvalidator.ValidateFieldsError, xvalidator.UtTranslator) (result map[string]string, need4xx bool) {
				return nil, false
			})}, nil, false},
		{"WithTranslatableError", translatableError("TODO"), []TranslateOption{}, map[string]string{"_": "TODO"}, true},
		{"WithTranslatableError", translatableError("TODO"), []TranslateOption{WithTranslatableError(
			func(e TranslatableError) (result map[string]string, need4xx bool) {
				return map[string]string{"_x_": e.Error()}, true
			})}, map[string]string{"_x_": "TODO"}, true},
		{"NilExtraErrors", errors.New("TODO"), []TranslateOption{WithExtraErrorsTranslate(nil)}, nil, false},
		{"ExtraErrors", errors.New("TODO"), []TranslateOption{WithExtraErrorsTranslate(
			func(e error) (result map[string]string, need4xx bool) {
				return map[string]string{"_": e.Error()}, true
			})}, map[string]string{"_": "TODO"}, true},
	} {
		t.Run("other_"+tc.name, func(t *testing.T) {
			result, need4xx := TranslateBindingError(tc.giveErr, tc.giveOpts...)
			xtesting.Equal(t, result, tc.wantMap)
			xtesting.Equal(t, need4xx, tc.want4xx)
		})
	}
}
