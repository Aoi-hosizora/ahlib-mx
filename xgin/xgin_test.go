package xgin

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xtesting"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"testing"
	"time"
)

func TestDumpRequest(t *testing.T) {
	app := gin.New()
	app.GET("a", func(c *gin.Context) {
		for _, s := range DumpRequest(c, nil) {
			log.Println(s)
		}
	})
	_ = app.Run(":1234")
}

func TestLogger(t *testing.T) {
	app := gin.New()

	l1 := logrus.New()
	l1.SetFormatter(&logrus.TextFormatter{ForceColors: true, TimestampFormat: time.RFC3339, FullTimestamp: true})
	l2 := log.New(os.Stderr, "", log.LstdFlags)

	PprofWrap(app)
	app.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		LogToLogrus(l1, c, start, end)
		LogToLogrus(l1, c, start, end, WithExtraText("abc"))
		LogToLogrus(l1, c, start, end, WithExtraFields(map[string]interface{}{"a": "b"}))
		LogToLogrus(l1, c, start, end, WithExtraFieldsV("a", "b"))
		LogToLogrus(l1, c, start, end, WithExtraText("abc"), WithExtraFieldsV("a", "b"))

		LogToLogger(l2, c, start, end)
		LogToLogger(l2, c, start, end, WithExtraText("abc"))
		LogToLogger(l2, c, start, end, WithExtraFields(map[string]interface{}{"a": "b"}))
	})
	app.GET("", func(c *gin.Context) {
		c.JSON(200, &gin.H{"ok": true})
	})

	_ = app.Run(":1234")
}

func TestGinLogger(t *testing.T) {
	app := gin.New()
	gin.ForceConsoleColor()
	app.Use(gin.Logger())
	app.GET("", func(c *gin.Context) {})
	app.GET("a", func(c *gin.Context) {})
	app.GET("a/:id", func(c *gin.Context) {})
	_ = app.Run(":1234")
	/*
		[GIN] 2021/01/26 - 12:03:53 | 200 |       956.9µs |             ::1 | GET      "/a/b"
		[GIN] 2021/01/26 - 12:04:28 | 404 |            0s |             ::1 | POST     "/a"
	*/
}

// func TestBinding(t *testing.T) {
// 	app := gin.New()
// 	enTrans, err := GetTranslator(en.New(), en_translations.RegisterDefaultTranslations)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
//
// 	log.Println(EnableRegexpBindingWithTranslator(enTrans))
// 	log.Println(EnableRFC3339DateBindingWithTranslator(enTrans))
// 	log.Println(EnableRFC3339DateTimeBindingWithTranslator(enTrans))
//
// 	type st struct {
// 		A string `binding:"regexp=^[abc]+$"`
// 		B string `binding:"date"`
// 		C string `binding:"datetime"`
// 		D string `binding:"gte=2"`
// 	}
//
// 	app.GET("", func(ctx *gin.Context) {
// 		log.Println(ctx.Request.RequestURI)
// 		st := &st{}
// 		err := ctx.ShouldBindQuery(st)
// 		if err != nil {
// 			translations := err.(validator.ValidationErrors).Translate(enTrans)
// 			ctx.JSON(200, &gin.H{
// 				"msg":    err.Error(),
// 				"detail": translations,
// 			})
// 		}
// 	})
//
// 	// http://localhost:1234/?A=a&B=2020-11-16&C=2020-11-16T21:44:03Z&D=555
// 	/*
// 		"detail": {
// 			"st.A": "A must matches regexp /^[abc]+$/",
// 			"st.B": "B must be an RFC3339 Date",
// 			"st.C": "C must be an RFC3339 DateTime",
// 			"st.D": "D must be at least 2 characters in length"
// 		},
// 	*/
// 	_ = app.Run(":1234")
// }

func TestAppRouter(t *testing.T) {
	app := gin.New()
	app.HandleMethodNotAllowed = true
	app.NoRoute(func(c *gin.Context) { c.String(200, "404 %s not found", c.Request.URL.String()) })
	app.NoMethod(func(c *gin.Context) { c.String(200, "405 %s not allowed", c.Request.Method) })

	g := app.Group("v1")
	ar := NewAppRouter(app, g)

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
	ar.Register()

	xtesting.Panic(t, func() {
		g := app.Group("v2/_test")
		ar := NewAppRouter(app, g)
		ar.GET("", func(c *gin.Context) {})
		ar.Register()
	})

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
