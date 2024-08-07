package trace

import (
	"bytes"
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
)

// Handler will set `trace_id` in context.
func Handler(c *gin.Context) {
	if _, ok := c.Get("trace_id"); !ok {
		ctx := context.WithValue(c.Request.Context(), "trace_id", util.GenerateRandomString(5, util.NUMBER_RUNES))
		c.Request = c.Request.WithContext(ctx)
	}
	c.Next()
}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)
		traceId := c.Request.Context().Value("trace_id")
		// headers
		log.Debugf("[REQUEST] trace_id: %v, headers: %v", traceId, c.Request.Header)
		// body
		log.Debugf("[REQUEST] trace_id: %v, body: %s", traceId, string(body))
		c.Next()
	}
}
