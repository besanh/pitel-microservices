package jwt_lib

import "github.com/golang-jwt/jwt/v5"

var SIGNED_KEY = []byte("pitel-pusher-secret-key")

func GenerateJWT(claims *jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret"))
}
