package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w BodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

		bodyLogWriter := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		start := time.Now()
		c.Next()
		end := time.Now()
		responseBody := bodyLogWriter.body.String()
		logField := map[string]interface{}{
			"uri":             c.Request.URL.Path,
			"raw_query":       c.Request.URL.RawQuery,
			"start_timestamp": start.Format("2006-01-02 15:04:05"),
			"end_timestamp":   end.Format("2006-01-02 15:04:05"),
			"server_name":     c.Request.Host,
			"remote_addr":     c.ClientIP(),
			"proto":           c.Request.Proto,
			"referer":         c.Request.Referer(),
			"request_method":  c.Request.Method,
			"response_time":   end.Sub(start).Milliseconds(),
			"content_type":    c.Request.Header.Get("Content-Type"),
			"status":          c.Writer.Status(),
			"user_agent":      c.Request.UserAgent(),
			//"request_body":    string(requestBody),
			"response_body": responseBody,
			"response_err":  c.Errors.Last(),
		}

		bf2 := bytes.NewBuffer([]byte{})
		jsonEncoder := json.NewEncoder(bf2)
		jsonEncoder.SetEscapeHTML(false)
		jsonEncoder.SetIndent("", "\t")
		jsonEncoder.Encode(logField)
		fmt.Println(bf2.String())
	}
}
