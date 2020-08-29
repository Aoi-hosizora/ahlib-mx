package xgin

import (
	"fmt"
	"github.com/Aoi-hosizora/ahlib/xnumber"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"math"
	"time"
)

func WithLogrus(logger *logrus.Logger, start time.Time, c *gin.Context) {
	latency := time.Now().Sub(start)
	method := c.Request.Method
	path := c.Request.URL.Path
	ip := c.ClientIP()
	code := c.Writer.Status()
	length := math.Abs(float64(c.Writer.Size()))
	lengthStr := xnumber.RenderByte(length)

	entry := logger.WithFields(logrus.Fields{
		"module":   "gin",
		"method":   method,
		"path":     path,
		"latency":  latency,
		"code":     code,
		"length":   length,
		"clientIP": ip,
	})

	if len(c.Errors) == 0 {
		msg := fmt.Sprintf("[Gin] %8d | %12s | %15s | %10s | %-7s %s", code, latency.String(), ip, lengthStr, method, path)
		if code >= 500 {
			entry.Error(msg)
		} else if code >= 400 {
			entry.Warn(msg)
		} else {
			entry.Info(msg)
		}
	} else {
		msg := fmt.Sprintf("[Gin] %s", c.Errors.ByType(gin.ErrorTypePrivate).String())
		entry.Error(msg)
	}
}

func WithLogger(logger *log.Logger, start time.Time, c *gin.Context) {
	latency := time.Now().Sub(start)
	method := c.Request.Method
	path := c.Request.URL.Path
	ip := c.ClientIP()
	code := c.Writer.Status()
	length := math.Abs(float64(c.Writer.Size()))
	latencyStr := latency.String()
	lengthStr := xnumber.RenderByte(length)

	if len(c.Errors) == 0 {
		msg := fmt.Sprintf("[Gin] %8d | %12s | %15s | %10s | %-7s %s", code, latencyStr, ip, lengthStr, method, path)
		logger.Println(msg)
	} else {
		msg := c.Errors.ByType(gin.ErrorTypePrivate).String()
		logger.Println(fmt.Sprintf("[Gin] %s", msg))
	}
}
