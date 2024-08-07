package responsetime

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/tel4vn/fins-microservices/common/hash"
	"github.com/tel4vn/fins-microservices/internal/fingerprint"
)

// Handler will set `X-Response-Time` header in response.
func Handler(c *gin.Context) {
	userAgent := c.Request.UserAgent()
	fp := fingerprint.FingerprintMD(c.Request)
	token := c.GetHeader("Authorization")
	if len(fp) < 1 {
		fp = "no fingerprint"
	}
	c.Set("fingerprint", hash.HashMD5(fp))
	c.Set("user_agent", userAgent)
	c.Set("token", token)
	c.Next()
}

func FingerprintMiddleware(ctx huma.Context, next func(huma.Context)) {
	fp := fingerprint.FingerprintHuma(ctx)
	token := ctx.Header("Authorization")
	if len(fp) < 1 {
		fp = "no fingerprint"
	}
	ctx = huma.WithValue(ctx, "fingerprint", hash.HashMD5(fp))
	ctx = huma.WithValue(ctx, "user_agent", ctx.Header("User-Agent"))
	ctx = huma.WithValue(ctx, "token", token)
	next(ctx)
}
