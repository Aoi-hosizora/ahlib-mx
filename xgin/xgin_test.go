package xgin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	logrus2 "github.com/sirupsen/logrus"
	"log"
	"os"
	"testing"
	"time"
)

func TestDumpRequest(t *testing.T) {
	app := gin.New()
	app.GET("a", func(c *gin.Context) {
		for _, s := range DumpRequest(c) {
			log.Println(s)
		}
	})
	_ = app.Run(":1234")
}

func TestBuildErrorDto(t *testing.T) {
	app := gin.New()
	app.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				e := BuildErrorDto(err, c, 2, true)
				e.Others = map[string]interface{}{"a": "b"}
				c.JSON(200, e)
			}
		}()
		c.Next()
	})
	app.GET("panic", func(c *gin.Context) {
		panic("test panic")
	})
	app.GET("error", func(c *gin.Context) {
		c.JSON(200, BuildBasicErrorDto(fmt.Errorf("test error"), c))
	})
	_ = app.Run(":1234")
}

func TestLogger(t *testing.T) {
	app := gin.New()

	logger := log.New(os.Stderr, "", log.LstdFlags)
	logrus := logrus2.New()

	PprofWrap(app)
	app.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()

		WithLogrus(logrus, start, c, nil)
		WithLogrus(logrus, start, c, WithExtraString("abc"))
		WithLogrus(logrus, start, c, WithExtraFields(map[string]interface{}{"a": "b"}))
		WithLogrus(logrus, start, c, WithExtraFieldsV("a", "b"))
		WithLogrus(logrus, start, c, WithExtraString("abc"), WithExtraFields(map[string]interface{}{"a": "b"}))

		WithLogger(logger, start, c, nil)
		WithLogger(logger, start, c, WithExtraString("abc"))
		WithLogger(logger, start, c, WithExtraFields(map[string]interface{}{"a": "b"}))
	})

	_ = app.Run(":1234")
}

func TestBinding(t *testing.T) {
	app := gin.New()
	enTrans, err := GetTranslator(en.New(), en_translations.RegisterDefaultTranslations)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(EnableRegexpBindingWithTranslator(enTrans))
	log.Println(EnableRFC3339DateBindingWithTranslator(enTrans))
	log.Println(EnableRFC3339DateTimeBindingWithTranslator(enTrans))

	type st struct {
		A string `binding:"regexp=^[abc]+$"`
		B string `binding:"date"`
		C string `binding:"datetime"`
		D string `binding:"gte=2"`
	}

	app.GET("", func(ctx *gin.Context) {
		log.Println(ctx.Request.RequestURI)
		st := &st{}
		err := ctx.ShouldBindQuery(st)
		if err != nil {
			translations := err.(validator.ValidationErrors).Translate(enTrans)
			ctx.JSON(200, &gin.H{
				"msg":    err.Error(),
				"detail": translations,
			})
		}
	})

	// http://localhost:1234/?A=a&B=2020-11-16&C=2020-11-16T21:44:03Z&D=555
	/*
		"detail": {
			"st.A": "A must matches regexp /^[abc]+$/",
			"st.B": "B must be an RFC3339 Date",
			"st.C": "C must be an RFC3339 DateTime",
			"st.D": "D must be at least 2 characters in length"
		},
	*/
	_ = app.Run(":1234")
}

func TestRoute(t *testing.T) {
	app := gin.New()
	app.HandleMethodNotAllowed = true
	app.NoRoute(func(c *gin.Context) { c.String(200, "404 %s not found", c.Request.URL.String()) })
	app.NoMethod(func(c *gin.Context) { c.String(200, "405 %s not allowed", c.Request.Method) })

	g := app.Group("v1")
	ar := NewAppRoute(app, g)

	ar.GET("", func(c *gin.Context) { log.Println(0, c.FullPath()) })
	ar.GET("a", func(c *gin.Context) { log.Println(1, c.FullPath()) })
	ar.GET("b", func(c *gin.Context) { log.Println(2, c.FullPath()) })
	ar.GET(":a", func(c *gin.Context) { log.Println(3, c.FullPath(), "|", c.Param("a")) })
	ar.GET("a/b", func(c *gin.Context) { log.Println(4, c.FullPath()) })
	ar.GET("c/d", func(c *gin.Context) { log.Println(5, c.FullPath()) })
	ar.GET(":a/:b", func(c *gin.Context) { log.Println(6, c.FullPath(), "|", c.Param("a"), "|", c.Param("b")) })
	ar.GET("a/b/c", func(c *gin.Context) { log.Println(7, c.FullPath()) })
	ar.GET("d/e/f", func(c *gin.Context) { log.Println(8, c.FullPath()) })
	ar.GET(":a/:b/:c", func(c *gin.Context) { log.Println(9, c.FullPath(), "|", c.Param("a"), "|", c.Param("b"), "|", c.Param("c")) })
	ar.GET("a/b/c/d", func(c *gin.Context) { log.Println(10, c.FullPath()) })
	ar.POST("a", func(c *gin.Context) { log.Println(11, c.FullPath()) }, func(c *gin.Context) { log.Println(11, c.FullPath()) })
	ar.POST("a/b", func(c *gin.Context) { log.Println(12, c.FullPath()) })
	ar.POST("a/b/:c", func(c *gin.Context) { log.Println(13, c.FullPath(), "|", c.Param("c")) })
	ar.Do()

	// TODO curl -X POST localhost:1234/v1/a/b/c/dd 405 POST not allowed

	_ = app.Run(":1234")
}

func TestRequiredAndOmitempty(t *testing.T) {
	unmarshal := func(obj interface{}, j string) interface{} {
		err := json.Unmarshal([]byte(j), obj)
		if err != nil {
			log.Fatalln(err)
		}
		return obj
	}
	type S1 struct {
		A uint64 `binding:"required"`
		B string `binding:"required"`
	}
	type S2 struct {
		A *uint64 `binding:"required"`
		B *string `binding:"required"`
	}
	type S3 struct {
		A uint64 `binding:"omitempty"`
		B string `binding:"omitempty"`
	}
	type S4 struct {
		A *uint64 `binding:"omitempty"`
		B *string `binding:"omitempty"`
	}
	type S5 struct {
		A uint64 `binding:"required,omitempty"`
		B string `binding:"required,omitempty"`
	}
	type S6 struct {
		A *uint64 `binding:"required,omitempty"`
		B *string `binding:"required,omitempty"`
	}

	v := validator.New()
	v.SetTagName("binding")

	// string required
	fmt.Println()
	log.Println(v.Struct(unmarshal(&S1{}, `{}`)) == nil)                     // false
	log.Println(v.Struct(unmarshal(&S1{}, `{"A": null, "B": null}`)) == nil) // false
	log.Println(v.Struct(unmarshal(&S1{}, `{"A": 0, "B": ""}`)) == nil)      // false
	log.Println(v.Struct(unmarshal(&S1{}, `{"A": 1, "B": " "}`)) == nil)     // true
	// *string required
	fmt.Println()
	log.Println(v.Struct(unmarshal(&S2{}, `{}`)) == nil)                     // false
	log.Println(v.Struct(unmarshal(&S2{}, `{"A": null, "B": null}`)) == nil) // false
	log.Println(v.Struct(unmarshal(&S2{}, `{"A": 0, "B": ""}`)) == nil)      // true
	log.Println(v.Struct(unmarshal(&S2{}, `{"A": 1, "B": " "}`)) == nil)     // true
	// string omitempty
	fmt.Println()
	log.Println(v.Struct(unmarshal(&S3{}, `{}`)) == nil)                     // true
	log.Println(v.Struct(unmarshal(&S3{}, `{"A": null, "B": null}`)) == nil) // true
	log.Println(v.Struct(unmarshal(&S3{}, `{"A": 0, "B": ""}`)) == nil)      // true
	log.Println(v.Struct(unmarshal(&S3{}, `{"A": 1, "B": " "}`)) == nil)     // true
	// *string omitempty => string omitempty
	fmt.Println()
	log.Println(v.Struct(unmarshal(&S4{}, `{}`)) == nil)                     // true
	log.Println(v.Struct(unmarshal(&S4{}, `{"A": null, "B": null}`)) == nil) // true
	log.Println(v.Struct(unmarshal(&S4{}, `{"A": 0, "B": ""}`)) == nil)      // true
	log.Println(v.Struct(unmarshal(&S4{}, `{"A": 1, "B": " "}`)) == nil)     // true
	// string required,omitempty => string required
	fmt.Println()
	log.Println(v.Struct(unmarshal(&S5{}, `{}`)) == nil)                     // false
	log.Println(v.Struct(unmarshal(&S5{}, `{"A": null, "B": null}`)) == nil) // false
	log.Println(v.Struct(unmarshal(&S5{}, `{"A": 0, "B": ""}`)) == nil)      // false
	log.Println(v.Struct(unmarshal(&S5{}, `{"A": 1, "B": " "}`)) == nil)     // true
	// *string required,omitempty => *string required
	fmt.Println()
	log.Println(v.Struct(unmarshal(&S6{}, `{}`)) == nil)                     // false
	log.Println(v.Struct(unmarshal(&S6{}, `{"A": null, "B": null}`)) == nil) // false
	log.Println(v.Struct(unmarshal(&S6{}, `{"A": 0, "B": ""}`)) == nil)      // true
	log.Println(v.Struct(unmarshal(&S6{}, `{"A": 1, "B": " "}`)) == nil)     // true
}
