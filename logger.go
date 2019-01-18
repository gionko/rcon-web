package main

import (
	"fmt"
	"time"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()

	log.Level = logrus.DebugLevel
	log.Out = os.Stdout
}

func LogMiddleware() gin.HandlerFunc {
	return LogMiddlewareWithLogger(logrus.StandardLogger())
}

func LogMiddlewareWithLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		ip := c.ClientIP()
		status := c.Writer.Status()
		latency := get_latency(start, time.Now())
		method := c.Request.Method
		comment := c.Errors.String()

		str := fmt.Sprintf("%15s | %10s | %3d | %-7s %s",
		                   ip,
		                   latency,
		                   status,
		                   method,
		                   path,
		)

		if len(comment) > 0 {
			str += fmt.Sprintf(" (%s)", comment)
		}

		logger.Info(str)
	}
}

func get_latency(start time.Time, end time.Time) string {
	var str string

	t := end.Sub(start)

	n := t.Nanoseconds()
	u := float64(n) / 1000
	m := u / 1000
	s := m / 1000

	if (n < 1000) {
		str = fmt.Sprintf("%7d ns", n)
	} else if (u < 1000) {
		str = fmt.Sprintf("%3.3f us", u)
	} else if (m < 1000) {
		str = fmt.Sprintf("%3.3f ms", m)
	} else {
		str = fmt.Sprintf("%3.3f s ", s)
	}

	return str
}
