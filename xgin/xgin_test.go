package xgin

import (
	"encoding/json"
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xstring"
	"github.com/gin-gonic/gin"
	"log"
	"testing"
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
				j, _ := json.Marshal(BuildErrorDto(err, c, 2, true))
				fmt.Println(xstring.PrettifyJson(string(j), 4, " "))
			}
		}()
		c.Next()
	})
	app.GET("", func(c *gin.Context) {
		panic("test panic")
	})
	_ = app.Run(":1234")
}

func TestW(t *testing.T) {
	app := gin.New()
	app.GET("json", JsonW(func(c *gin.Context) (int, interface{}) {
		return 200, &gin.H{"test": "hello world"}
	}))
	app.GET("xml", XmlW(func(c *gin.Context) (int, interface{}) {
		return 200, &gin.H{"test": "hello world"}
	}))
	_ = app.Run(":1234")
}
