package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/codestagea/bindmgr/internal/tools"
	"runtime/debug"
	"strconv"
	"strings"
)

func CustomError(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if c.IsAborted() {
				c.Status(200)
			}
			switch errStr := err.(type) {
			case string:
				p := strings.Split(errStr, "#")
				if len(p) == 3 && p[0] == "CustomError" {
					statusCode, e := strconv.Atoi(p[1])
					if e != nil {
						statusCode = 400
					}
					c.Status(statusCode)
					logrus.Errorf("[%s][%s %v] %v %s: %s",
						c.Request.Method,
						c.Request.URL,
						statusCode,
						c.Request.RequestURI,
						c.ClientIP(),
						p[2],
					)
					tools.RetOfErrMsg(c, statusCode, p[2])
				}
			default:
				logrus.Errorf("unknown error %v,\n%s", err, debug.Stack())
				tools.RetOfErrMsg(c, 400, "system error")
			}
		}
	}()
	c.Next()
}
