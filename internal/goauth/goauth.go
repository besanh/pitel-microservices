package goauth

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tel4vn/pitel-microservices/common/util"
)

type (
	IAuth interface {
		Add(ctx context.Context, data *AuthUser) (user *AuthUser, err error)
		Find(ctx context.Context, token string) (user *AuthUser, err error)
		RefreshToken(ctx context.Context, refreshToken, token string) (user *AuthUser, err error)
		DeleteFullWithId(ctx context.Context, id string) (err error)
	}

	Auth struct {
		redisClient *redis.Client
	}

	AuthUser struct {
		Id           string    `json:"id"`
		UserId       string    `json:"user_id"`
		Fingerprint  string    `json:"fingerprint"`
		UserAgent    string    `json:"user_agent"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		CreatedTime  time.Time `json:"created_time"`
		ExpiredTime  time.Time `json:"expired_time"`
		Data         any       `json:"data"`
		TokenType    string    `json:"token_type"`
		Level        string    `json:"level"`
	}

	RefreshTokenValue struct {
		Id          string    `json:"id"`
		UserId      string    `json:"user_id"`
		Token       string    `json:"token"`
		ExpiredTime time.Time `json:"expired_time"`
	}
)

func (user *AuthUser) InitToken() {
	currentTime := time.Now().Local()
	expiredTime := currentTime.Add(time.Duration(TOKEN_TTL) * time.Second)
	user.Token = GenerateToken(user.Id)
	user.RefreshToken = GenerateRefreshToken(user.Id)
	user.ExpiredTime = expiredTime
	user.CreatedTime = currentTime
	user.TokenType = TOKEN_TYPE
}

const (
	TOKEN_KEY         = "pitel_bss_chat_auth_token"
	USER_KEY          = "pitel_bss_chat_auth_user"
	REFRESH_TOKEN_KEY = "pitel_bss_chat_auth_refresh_token"
	TOKEN_TTL         = 86400
	// TOKEN_TTL         = 120
	REFRESH_TOKEN_TTL = 172800
	TOKEN_TYPE        = "Bearer"
)

func NewGoAuth(client *redis.Client) IAuth {
	return &Auth{
		redisClient: client,
	}
}

var GoAuthClient IAuth

func (g *Auth) getAuthUserByToken(ctx context.Context, token string) (user *AuthUser, err error) {
	result, err := g.redisClient.HMGet(ctx, TOKEN_KEY, token).Result()
	if err != nil {
		return
	} else if len(result) < 1 {
		return nil, nil
	}
	user = new(AuthUser)
	val, ok := result[0].(string)
	if !ok {
		return nil, nil
	} else if err = util.ParseStringToAny(val, user); err != nil {
		return
	}
	return
}

func (g *Auth) getAuthUserById(ctx context.Context, userId string) (user *AuthUser, err error) {
	result, err := g.redisClient.HMGet(ctx, USER_KEY, userId).Result()
	if err != nil {
		return
	} else if len(result) < 1 {
		return nil, nil
	}
	user = new(AuthUser)
	val, ok := result[0].(string)
	if !ok {
		return nil, nil
	} else if err = util.ParseStringToAny(val, user); err != nil {
		return
	}
	return
}

func (g *Auth) getValueOfRefreshToken(ctx context.Context, refreshToken string) (value *RefreshTokenValue, err error) {
	result, err := g.redisClient.HMGet(ctx, REFRESH_TOKEN_KEY, refreshToken).Result()
	if err != nil {
		return
	} else if len(result) < 1 {
		return nil, nil
	}
	value = new(RefreshTokenValue)
	val, ok := result[0].(string)
	if !ok {
		return nil, nil
	} else if err = util.ParseStringToAny(val, value); err != nil {
		return
	}
	return
}

func (g *Auth) Find(ctx context.Context, token string) (user *AuthUser, err error) {
	user, err = g.getAuthUserByToken(ctx, token)
	if err != nil {
		return
	} else if user == nil {
		return
	}
	currentTime := time.Now().Local()
	if user.ExpiredTime.Sub(currentTime) <= 0 {
		err = errors.New("token is expired")
		return
	}
	return
}

func (g *Auth) Add(ctx context.Context, data *AuthUser) (user *AuthUser, err error) {
	user, err = g.getAuthUserById(ctx, data.Id)
	if err != nil {
		return
	} else if user == nil {
		data.InitToken()
		user = data
		if err = g.addUserAuth(ctx, *user); err != nil {
			return
		}
	} else if user.ExpiredTime.Sub(time.Now().Local()) <= 0 {
		// token is expired
		if err = g.deleteAuthUser(ctx, *user); err != nil {
			return
		}
		data.InitToken()
		user = data
		if err = g.addUserAuth(ctx, *user); err != nil {
			return
		}
	}
	return
}

func (g *Auth) addUserAuth(ctx context.Context, data AuthUser) (err error) {
	value, err := util.ParseAnyToString(data)
	if err != nil {
		return
	}
	if err = g.redisClient.HSet(ctx, USER_KEY, data.Id, value).Err(); err != nil {
		return
	}
	if err = g.redisClient.HSet(ctx, TOKEN_KEY, data.Token, value).Err(); err != nil {
		return
	}
	refreshTokenValue := RefreshTokenValue{
		UserId:      data.Id,
		Token:       data.RefreshToken,
		ExpiredTime: time.Now().Local().Add(time.Duration(REFRESH_TOKEN_TTL) * time.Second),
	}
	value, err = util.ParseAnyToString(refreshTokenValue)
	if err != nil {
		return
	}
	if err = g.redisClient.HSet(ctx, REFRESH_TOKEN_KEY, data.RefreshToken, value).Err(); err != nil {
		return
	}
	return
}

func (g *Auth) deleteAuthUser(ctx context.Context, user AuthUser) (err error) {
	if err = g.redisClient.HDel(ctx, USER_KEY, user.Id).Err(); err != nil {
		return
	}
	if err = g.redisClient.HDel(ctx, TOKEN_KEY, user.Token).Err(); err != nil {
		return
	}
	if err = g.deleteRefreshToken(ctx, user.RefreshToken); err != nil {
		return
	}
	return
}

func (g *Auth) deleteRefreshToken(ctx context.Context, refreshToken string) (err error) {
	if err = g.redisClient.HDel(ctx, REFRESH_TOKEN_KEY, refreshToken).Err(); err != nil {
		return
	}
	return
}

func (g *Auth) RefreshToken(ctx context.Context, refreshToken, token string) (user *AuthUser, err error) {
	refreshTokenValue, err := g.getValueOfRefreshToken(ctx, refreshToken)
	if err != nil {
		return
	} else if refreshTokenValue == nil {
		return
	}
	currentTime := time.Now().Local()
	if refreshTokenValue.ExpiredTime.Sub(currentTime) <= 0 {
		err = errors.New("refresh token is expired")
		return
	}
	user, err = g.getAuthUserById(ctx, refreshTokenValue.UserId)
	if err != nil {
		return
	} else if user == nil {
		return
	}
	// if user.RefreshToken != refreshToken {
	// 	if err = g.deleteRefreshToken(ctx, refreshToken); err != nil {
	// 		return
	// 	}
	// 	err = errors.New("refresh token is invalid")
	// 	return
	// }
	// currentToken := user.Token
	user.InitToken()
	// Not renew old token, just renew refresh token
	// user.Token = currentToken
	if err = g.addUserAuth(ctx, *user); err != nil {
		return
	}
	if err = g.deleteRefreshToken(ctx, refreshToken); err != nil {
		return
	}
	return

}

func GenerateToken(id string) string {
	uuidNew, _ := uuid.NewRandom()
	idEnc := base64.StdEncoding.EncodeToString([]byte(id))
	token := strings.Replace(uuidNew.String(), "-", "", -1)
	token = token + "-" + idEnc
	return token
}

func GenerateRefreshToken(id string) string {
	uuidNew, _ := uuid.NewRandom()
	idEnc := base64.StdEncoding.EncodeToString([]byte(id))
	token := strings.Replace(uuidNew.String(), "-", "", -1)
	token = "fre-" + token + idEnc
	return token
}

func (g *Auth) DeleteFullWithId(ctx context.Context, id string) (err error) {
	user, err := g.getAuthUserById(ctx, id)
	if err != nil {
		return
	} else if user == nil {
		return
	}
	if err = g.deleteAuthUser(ctx, *user); err != nil {
		return
	}
	return
}
