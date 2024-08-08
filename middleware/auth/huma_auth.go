package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
)

var DefaultAuthSecurity = []map[string][]string{
	{"ibkAuth": {""}},
}

func NewAuthMiddleware(api hureg.APIGen) func(ctx huma.Context, next func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		isAuthorizationRequired := false
		for _, opScheme := range ctx.Operation().Security {
			var ok bool
			if _, ok = opScheme["ibkAuth"]; ok {
				isAuthorizationRequired = true
				break
			}
		}
		log.DebugfContext(ctx.Context(), "is require authorization: %v", isAuthorizationRequired)
		if isAuthorizationRequired {
			HumaAuthMiddleware(api, ctx, next)
		} else {
			next(ctx)
		}
	}
}

func HumaAuthMiddleware(api hureg.APIGen, ctx huma.Context, next func(huma.Context)) {
	var username, password string

	headerValue := ctx.Header("Authorization")
	basicAuthPrefix := "Basic "
	if strings.HasPrefix(headerValue, basicAuthPrefix) {
		encodedCreds := headerValue[len(basicAuthPrefix):]
		creds, err := base64.StdEncoding.DecodeString(encodedCreds)
		if err != nil {
			huma.WriteErr(api.GetHumaAPI(), ctx, http.StatusUnauthorized, string(constants.ERR_UNAUTHORIZED), err)
			return
		}
		credsParts := strings.SplitN(string(creds), ":", 2)
		if len(credsParts) < 2 {
			huma.WriteErr(api.GetHumaAPI(), ctx, http.StatusUnauthorized, string(constants.ERR_UNAUTHORIZED), err)
			return
		}

		username, password = credsParts[0], credsParts[1]
		// TODO: implement auth based on token
		_ = username
		_ = password
	} else {
		token := parseTokenFromAuthorization(headerValue)
		user, err := validateToken(ctx.Context(), token)
		if err != nil {
			huma.WriteErr(api.GetHumaAPI(), ctx, http.StatusUnauthorized, string(constants.ERR_UNAUTHORIZED), err)
			return
		}
		ctx = huma.WithValue(ctx, "user", user)
	}

	next(ctx)
}

func parseTokenFromAuthorization(authorizationHeader string) string {
	return strings.Replace(authorizationHeader, "Bearer ", "", 1)
}
