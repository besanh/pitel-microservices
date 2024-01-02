package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SIGNATURE = "d6wSGXochuK9v5V9dDPch1hsSeY0xpiMgHVJkATRsdjgnpUasG"
)

func ValidHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("ICOMNI-Signature")
		if len(header) > 0 {
		} else if len(ctx.GetHeader("Authorization")) > 0 {
			header = ctx.GetHeader("Authorization")
		}
		token := header
		if token != SIGNATURE {
			ctx.JSON(
				http.StatusUnauthorized,
				map[string]any{
					"error": http.StatusText(http.StatusUnauthorized),
				},
			)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
