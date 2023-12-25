package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidHeader(incomSignature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("ICOMNI-Signature")
		if token != incomSignature {
			c.JSON(
				http.StatusUnauthorized,
				map[string]interface{}{
					"error": http.StatusText(http.StatusUnauthorized),
				},
			)
			c.Abort()
			return
		}
		c.Next()
	}
}
