package logger

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type bufferedWriter struct {
	gin.ResponseWriter
	out    *bufio.Writer
	Buffer bytes.Buffer
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// set gin log
func InitGinLogger() gin.HandlerFunc {
	// 不需要记录日志的路由，如静态资源文件
	notLoggerPath := []string{
		"/static",
	}
	return func(c *gin.Context) {

		for _, v := range notLoggerPath {
			if strings.HasPrefix(c.Request.URL.Path, v) {
				c.Next()
				return
			}
		}

		xTraceId := c.Request.Header.Get("req-id")
		if xTraceId == "" {
			xTraceId = DefaultGenRequestID()
		}
		c.Header("req-id", xTraceId)
		c.Request.Header.Set("req-id", xTraceId)

		Std = New(xTraceId).Caller(4)

		// 记录请求日志
		guest := c.Request.Header.Get("Authentication")
		params, _ := io.ReadAll(c.Request.Body)
		reqBytes, _ := httputil.DumpRequest(c.Request, true)
		reqLoggerFields := map[string]interface{}{
			"guest": guest, // 身份标识
			"ip":    c.ClientIP(),
		}

		WithFields(reqLoggerFields).Caller(3).Info(fmt.Sprintf(
			"%s \n ****** params ******* \n %s\n ",
			string(reqBytes),
			string(params[:300]),
		))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(params))
		if c.Writer.Status() == 200 {
			w := bufio.NewWriter(c.Writer)
			buff := bytes.Buffer{}
			newWriter := &bufferedWriter{c.Writer, w, buff}
			c.Writer = newWriter

			defer func() {
				response := newWriter.Buffer.Bytes()
				reqLoggerFields := map[string]interface{}{
					"Status":   c.Writer.Status(),
					"guest":    guest, // 身份标识
					"clientIP": c.ClientIP(),
				}
				if len(response) > 0 {
					reqLoggerFields["Response"] = response
				}
				WithFields(reqLoggerFields).Caller(4).Info(string(reqBytes))
				w.Flush()
			}()
		}

		t := time.Now()

		c.Next()
		username, _ := c.Get("username")
		// 请求后
		latency := time.Since(t)
		statusCode := c.Writer.Status()
		statusColor := colorForStatus(statusCode)
		method := c.Request.Method
		methodColor := colorForMethod(method)
		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		resLoggerFields := map[string]interface{}{
			"guest":    guest, // 身份标识
			"clientIP": c.ClientIP(),
			"comment":  comment,
			"username": username,
		}
		WithFields(resLoggerFields).Caller(3).Info(fmt.Sprintf(
			"|%s %s %s| |%s %3d %s||%s|",
			methodColor, method, reset,
			statusColor, statusCode, reset,
			latency,
		))
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}
