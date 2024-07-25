package auth

// func (mdw *GatewayAuthMiddleware) HTTPAuthMiddleware() gin.HandlerFunc {
// 	return mdw.HTTPAuthMiddlewareContext
// }

// func (mdw *GatewayAuthMiddleware) HTTPAuthMiddlewareContext(c *gin.Context) {
// 	isAuthenticaed, ok := c.Value(AUTHENTICATED).(bool)
// 	if !ok || !isAuthenticaed {
// 		md := HandleMetadata(c, c.Request)
// 		ctx := metadata.NewIncomingContext(c, md)
// 		c.Request = c.Request.WithContext(ctx)
// 		var user *model.AuthUser
// 		if mdw.env == DEV {
// 			user = ParseHeaderToUserDev(c.Request.Context())
// 		} else {
// 			user = ParseHeaderToUser(c.Request.Context())
// 		}
// 		if user == nil || len(user.UserId) < 1 {
// 			c.JSON(response.Unauthorized())
// 			c.Abort()
// 			return
// 		}
// 		c.Set(AUTHENTICATED, true)
// 		c.Set(USER, user)
// 	}
// }
