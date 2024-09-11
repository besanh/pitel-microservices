package service

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tel4vn/pitel-microservices/common/cache"
	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/response"
	"github.com/tel4vn/pitel-microservices/internal/goauth"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/repository"
)

type (
	IChatAuth interface {
		Login(ctx context.Context, body model.LoginRequest) (result *model.LoginResponse, err error)
		RefreshToken(ctx context.Context, body model.RefreshTokenRequest) (result *model.RefreshTokenResponse, err error)
		VerifyToken(ctx context.Context, token string) (goAuthUser *goauth.AuthUser, err error)
		Logout(ctx context.Context, token, fingerprint string) (err error)
	}
	ChatAuth struct{}
)

const (
	STATE_REQUEST      = "request"
	STATE_VERIFY_TOKEN = "verify_token"
	STATE_VERIFY_CODE  = "verify_code"
	STATE_SUBMIT       = "submit"
)

var ChatAuthService IChatAuth

func NewChatAuth() IChatAuth {
	return &ChatAuth{}
}

func (auth *ChatAuth) Login(ctx context.Context, body model.LoginRequest) (result *model.LoginResponse, err error) {
	filter := model.ChatUserFilter{
		Username: body.Username,
	}
	user, err := repository.ChatUserRepo.GetChatUserByFilter(ctx, repository.DBConn, filter)
	if err != nil {
		log.Error(err)
		return nil, err
	} else if user == nil {
		return nil, errors.New(response.ERR_GET_FAILED)
	} else if !validatePassword(user.Salt, user.Password, body.Password) {
		return nil, errors.New(response.ERR_INVALID_USERNAME_PASSWORD)
	}

	// TODO: check tenant_id
	tokenData := model.TokenData{
		// TenantId: tenant.Id,
		Username: user.Username,
		UserId:   user.Id,
		RoleId:   user.RoleId,
		Level:    user.Level,
	}
	id := fmt.Sprintf("%s.%s", user.Id, body.Fingerprint)
	goAuthUser, err := goauth.GoAuthClient.Add(ctx, &goauth.AuthUser{
		Id:          id,
		Fingerprint: body.Fingerprint,
		UserAgent:   body.UserAgent,
		UserId:      user.Id,
		Data:        tokenData,
		Level:       user.Level,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result = &model.LoginResponse{
		Token:        goAuthUser.Token,
		RefreshToken: goAuthUser.RefreshToken,
		ExpiredIn:    int(time.Until(goAuthUser.ExpiredTime).Seconds()),
		Fullname:     user.Fullname,
		Username:     user.Username,
		UserId:       user.Id,
		RoleId:       user.RoleId,
		Level:        user.Level,
	}
	return result, nil
}

func validatePassword(salt, userPassword, password string) (ok bool) {
	ok = false
	arr := strings.Split(userPassword, "$")
	if len(arr) < 1 {
		return
	}
	hashPassword := arr[0]
	dbPassword := []byte(salt + password)
	encryptPassword := fmt.Sprintf("%x", md5.Sum(dbPassword))
	if hashPassword != encryptPassword {
		return
	}
	return true
}

func (auth *ChatAuth) RefreshToken(ctx context.Context, body model.RefreshTokenRequest) (result *model.RefreshTokenResponse, err error) {
	goAuthUser, err := goauth.GoAuthClient.RefreshToken(ctx, body.RefreshToken, body.Token)
	if err != nil {
		log.Error(err)
		return nil, err
	} else if goAuthUser == nil {
		return nil, errors.New(response.ERR_GET_FAILED)
	}
	result = &model.RefreshTokenResponse{
		RefreshToken: goAuthUser.RefreshToken,
		ExpiredIn:    int(time.Until(goAuthUser.ExpiredTime).Seconds()),
		Token:        goAuthUser.Token,
	}
	return result, nil
}

func (s *ChatAuth) VerifyToken(ctx context.Context, token string) (goAuthUser *goauth.AuthUser, err error) {
	if strings.Contains(token, "Basic ") {
		basicToken := strings.Replace(token, "Basic ", "", -1)
		cacheToken := cache.MCache.Get(basicToken)
		if cacheToken != nil {
			token, _ = cacheToken.(string)
		} else {
			bytes, _ := base64.StdEncoding.DecodeString(basicToken)
			arr := strings.Split(string(bytes), ":")
			if len(arr) < 2 {
				return nil, errors.New("token is invalid")
			}
			username := arr[0]
			password := arr[1]
			loginResponse, err := s.Login(ctx, model.LoginRequest{
				Username: username,
				Password: password,
			})
			if err != nil {
				log.Error(err)
				return nil, err
			}
			cache.MCache.Set(basicToken, loginResponse.Token, time.Hour)
			token = loginResponse.Token
		}
	}
	goAuthUser, err = goauth.GoAuthClient.Find(ctx, token)
	if err != nil {
		log.Error(err)
		return nil, err
	} else if goAuthUser == nil {
		return nil, errors.New("token is invalid")
	}
	return
}

// func (auth *ChatAuth) ForgotPassword(ctx context.Context, body model.ForgotPassword) (res *model.ForgotPasswordResponse, err error) {
// 	switch body.State {
// 	case STATE_REQUEST:
// 		res, err = auth.forgotPasswordRequest(ctx, body)
// 	case STATE_VERIFY_TOKEN, STATE_VERIFY_CODE, STATE_SUBMIT:
// 		res, err = auth.forgotPasswordSubmit(ctx, body)
// 	default:
// 		err = errors.New(response.ERR_REQUEST_INVALID)
// 	}
// 	return
// }

// func generateToken(id string) string {
// 	code := util.GenerateRandomString(36)
// 	idEnc := base64.StdEncoding.EncodeToString([]byte(id))
// 	token := fmt.Sprintf("chat-%s-%s", code, idEnc)
// 	return token
// }

// func (s *ChatAuth) forgotPasswordRequest(ctx context.Context, body model.ForgotPassword) (res *model.ForgotPasswordResponse, err error) {
// 	res = &model.ForgotPasswordResponse{
// 		IsValid: true,
// 	}
// 	param := model.Param{
// 		Key:      "email",
// 		Operator: "=",
// 		Value:    body.Email,
// 	}
// 	user, err := repository.UserRepo.FindByQuery(ctx, param)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	} else if user == nil {
// 		err = errors.New(response.ERR_EMAIL_NOTFOUND)
// 		return
// 	} else if len(user.Email) < 1 {
// 		err = errors.New(response.ERR_EMAIL_NOTFOUND)
// 		return
// 	}

// 	param = model.Param{
// 		Key:      "user_id",
// 		Operator: "=",
// 		Value:    user.Id,
// 	}
// 	current := time.Now()

// 	aaaRequest, err := repository.AAARequestRepo.FindByQuery(ctx, param)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	} else if aaaRequest != nil {
// 		isRemainTime := aaaRequest.CreatedAt.Add(1 * time.Minute).After(current)
// 		if isRemainTime {
// 			res = &model.ForgotPasswordResponse{
// 				IsValid:         false,
// 				Code:            "ERR_REQUEST_IS_EXISTED",
// 				NextRequestTime: aaaRequest.CreatedAt.Add(1 * time.Minute),
// 				// RemainTime:      current.Sub(aaaRequest.CreatedAt.Add(1 * time.Minute)),
// 			}
// 			return
// 		}
// 	}

// 	token := generateToken(user.Id)
// 	code := util.GenerateRandomString(6, util.NUMBER_RUNES)
// 	aaaRequest = &model.AAARequest{
// 		Base:         model.InitBaseModel("", ""),
// 		UserId:       user.Id,
// 		Method:       "forgot",
// 		Ip:           "",
// 		UserAgent:    body.UserAgent,
// 		Code:         code,
// 		Token:        token,
// 		InvalidCount: 0,
// 		ExpiredAt:    current.Add(30 * time.Minute),
// 		IsActive:     true,
// 	}

// 	if err = repository.AAARequestRepo.Insert(ctx, *aaaRequest); err != nil {
// 		log.Error(err)
// 		return
// 	}

// 	link := fmt.Sprintf("%s/forgot?token=%s", s.UIBaseUrl, token)
// 	log.Infof("link: %s", link)
// 	var emailBody string
// 	emailBody, err = mail.ParseTemplate("public/mail/forgot.html", struct {
// 		Link     string `json:"link"`
// 		Fullname string `json:"fullname"`
// 		Code     string `json:"code"`
// 	}{
// 		Link:     link,
// 		Fullname: user.Fullname,
// 		Code:     code,
// 	})
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// 	go func(email string, message string) {
// 		m := gomail.NewMessage()
// 		m.SetBody("text/html", message)
// 		if err = mail.Client.Send("[FinS] Forgot Password", email, m); err != nil {
// 			log.Error(err)
// 		}
// 	}(user.Email, emailBody)
// 	return
// }

// func (s *ChatAuth) forgotPasswordVerifyToken(ctx context.Context, body model.ForgotPassword) (res *model.ForgotPasswordResponse, err error) {
// 	param := model.Param{
// 		Key:      "token",
// 		Operator: "=",
// 		Value:    body.Token,
// 	}
// 	aaaRequest, err := repository.AAARequestRepo.FindByQuery(ctx, param)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	} else if aaaRequest == nil {
// 		err = errors.New(response.ERR_REQUEST_INVALID)
// 		return
// 	}

// 	if aaaRequest.ExpiredAt.Before(time.Now()) {
// 		err = errors.New(response.ERR_REQUEST_IS_EXPIRED)
// 		return
// 	}

// 	res = &model.ForgotPasswordResponse{
// 		IsValid: true,
// 	}

// 	return
// }

// func (s *ChatAuth) forgotPasswordSubmit(ctx context.Context, body model.ForgotPassword) (res *model.ForgotPasswordResponse, err error) {
// 	params := []model.Param{
// 		{
// 			Key:      "token",
// 			Operator: "=",
// 			Value:    body.Token,
// 		}, {
// 			Key:      "is_active",
// 			Operator: "=",
// 			Value:    true,
// 		},
// 	}
// 	aaaRequest, err := repository.AAARequestRepo.FindByQuery(ctx, params...)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	} else if aaaRequest == nil {
// 		err = errors.New(response.ERR_REQUEST_INVALID)
// 		return
// 	}
// 	res = &model.ForgotPasswordResponse{
// 		IsValid: true,
// 	}
// 	if aaaRequest.ExpiredAt.Before(time.Now()) {
// 		err = errors.New(response.ERR_REQUEST_IS_EXPIRED)
// 		return
// 	} else if body.State == STATE_VERIFY_TOKEN {
// 		return
// 	} else if body.Code != aaaRequest.Code {
// 		aaaRequest.InvalidCount = aaaRequest.InvalidCount + 1
// 		if err = repository.AAARequestRepo.Update(ctx, *aaaRequest); err != nil {
// 			log.Error(err)
// 		}
// 		err = errors.New(response.ERR_REQUEST_CODE_IS_INVALID)
// 		return
// 	} else if body.State == STATE_VERIFY_CODE {
// 		return
// 	}
// 	salt := util.GenerateRandomString(8, nil)
// 	password := fmt.Sprintf("%s$%s", salt, hashSalt(salt, body.Password))
// 	user, err := repository.UserRepo.FindById(ctx, aaaRequest.UserId)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	} else if user == nil {
// 		err = errors.New(response.ERR_REQUEST_INVALID)
// 		return
// 	}
// 	user.Password = password
// 	if err = repository.UserRepo.Update(ctx, *user); err != nil {
// 		log.Error(err)
// 		return
// 	}
// 	aaaRequest.IsActive = false
// 	if err = repository.AAARequestRepo.Update(ctx, *aaaRequest); err != nil {
// 		log.Error(err)
// 	}
// 	return
// }

func (s *ChatAuth) Logout(ctx context.Context, token, fingerprint string) (err error) {
	var goAuthUser *goauth.AuthUser
	goAuthUser, err = goauth.GoAuthClient.Find(ctx, token)
	if err != nil {
		log.Error(err)
		return err
	} else if goAuthUser == nil {
		return errors.New("token is invalid")
	}
	// if goAuthUser.Fingerprint != fingerprint {
	// 	return errors.New("token is invalid! mismatch fingerprint")
	// }
	err = goauth.GoAuthClient.DeleteFullWithId(ctx, goAuthUser.Id)
	if err != nil {
		log.Error(err)
		return err
	}
	return
}
