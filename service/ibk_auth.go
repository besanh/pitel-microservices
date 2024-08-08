package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/hash"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/internal/goauth"
	"github.com/tel4vn/fins-microservices/internal/mail"
	"github.com/tel4vn/fins-microservices/model"
	repository "github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/repository/sql_builder"
	"gopkg.in/gomail.v2"
)

type (
	IAuth interface {
		Login(ctx context.Context, body model.LoginRequest) (result *model.LoginResponse, err error)
		Logout(ctx context.Context, token, fingerprint string) (err error)
		RefreshToken(ctx context.Context, body model.RefreshTokenRequest) (result *model.RefreshTokenResponse, err error)
		VerifyToken(ctx context.Context, token string) (goAuthUser *goauth.AuthUser, err error)
		ForgotPassword(ctx context.Context, body model.ForgotPassword) (res *model.ForgotPasswordResponse, err error)
		RefreshAuthData(ctx context.Context, authUser *model.AuthUser, token string) (err error)
	}
	Auth struct {
		UIBaseUrl string
	}
)

var AuthService IAuth

func NewAuth(uiBaseUrl string) IAuth {
	return &Auth{
		UIBaseUrl: uiBaseUrl,
	}
}

func (auth *Auth) Login(ctx context.Context, body model.LoginRequest) (result *model.LoginResponse, err error) {
	user, err := repository.IBKUserRepo.GetUserByUsername(ctx, body.Username)
	if err != nil {
		log.Error(err)
		return nil, err
	} else if user == nil {
		return nil, variables.NewError(constants.ERR_USERNAME_NOT_FOUND)
	} else if !validatePassword(user.Password, body.Password) {
		return nil, variables.NewError(constants.ERR_INVALID_USERNAME_PASSWORD)
	}
	// group, err := repository.AAA_GroupRepo.FindById(ctx, user.GroupId)
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, err
	// } else if group == nil {
	// 	return nil, variables.NewError(constants.ERR_GROUP_NOTFOUND)
	// }
	// bu, err := repository.AAA_BusinessUnitRepo.FindById(ctx, user.BusinessUnitId)
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, err
	// } else if bu == nil {
	// 	return nil, variables.NewError(constants.ERR_BUSINESS_UNIT_NOTFOUND)
	// }
	// tenant, err := repository.AAA_TenantRepo.FindById(ctx, bu.TenantId)
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, err
	// } else if tenant == nil {
	// 	return nil, variables.NewError(constants.ERR_TENANT_NOTFOUND)
	// }
	// tokenData := model.AuthUserData{
	// 	TenantId:       tenant.Id,
	// 	BusinessUnitId: bu.Id,
	// 	Username:       user.Username,
	// 	UserId:         user.Id,
	// 	Level:          user.Level,
	// 	Scopes:         user.Scopes,
	// }
	// id := fmt.Sprintf("%s.%s", user.Id, body.Fingerprint)
	// goAuthUser, err := goauth.GoAuthClient.Add(ctx, &goauth.AuthUser{
	// 	Id:          id,
	// 	Fingerprint: body.Fingerprint,
	// 	UserAgent:   body.UserAgent,
	// 	UserId:      user.Id,
	// 	Data:        tokenData,
	// 	Level:       user.Level,
	// })
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, err
	// }
	// result = &model.LoginResponse{
	// 	Token:        goAuthUser.Token,
	// 	RefreshToken: goAuthUser.RefreshToken,
	// 	ExpiredIn:    int(time.Until(goAuthUser.ExpiredTime).Seconds()),
	// 	Fullname:     user.Fullname,
	// 	TenantName:   tenant.TenantName,
	// 	AuthUserData: &tokenData,
	// }
	return
}

func validatePassword(userPassword, password string) (ok bool) {
	ok = false
	arr := strings.Split(userPassword, "$")
	if len(arr) < 2 {
		return
	}
	salt := arr[0]
	hashPassword := arr[1]
	if hashPassword != hashSalt(salt, password) {
		return
	}
	return true
}

func hashSalt(salt string, str string) string {
	return hash.HashMD5(fmt.Sprintf("%s$%s", salt, str))
}

func (auth *Auth) RefreshToken(ctx context.Context, body model.RefreshTokenRequest) (result *model.RefreshTokenResponse, err error) {
	goAuthUser, err := goauth.GoAuthClient.RefreshToken(ctx, body.RefreshToken, body.Token)
	if err != nil {
		log.Error(err)
		return
	} else if goAuthUser == nil {
		err = variables.NewError(constants.ERR_REFRESH_TOKEN_INVALID)
		return
	}
	result = &model.RefreshTokenResponse{
		RefreshToken: goAuthUser.RefreshToken,
		ExpiredIn:    int(time.Until(goAuthUser.ExpiredTime).Seconds()),
		Token:        goAuthUser.Token,
	}
	return
}

func (auth *Auth) VerifyToken(ctx context.Context, token string) (goAuthUser *goauth.AuthUser, err error) {
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
			loginResponse, errTmp := auth.Login(ctx, model.LoginRequest{
				Username: username,
				Password: password,
			})
			if errTmp != nil {
				err = errTmp
				log.Error(err)
				return
			}
			cache.MCache.Set(basicToken, loginResponse.Token, time.Hour)
			token = loginResponse.Token
		}
	}
	goAuthUser, err = goauth.GoAuthClient.Find(ctx, token)
	if err != nil {
		log.Error(err)
		return
	} else if goAuthUser == nil {
		err = errors.New("token is invalid")
	}
	return
}

func (auth *Auth) ForgotPassword(ctx context.Context, body model.ForgotPassword) (res *model.ForgotPasswordResponse, err error) {
	switch body.State {
	case constants.STATE_REQUEST:
		res, err = auth.forgotPasswordRequest(ctx, body)
	case constants.STATE_VERIFY_TOKEN, constants.STATE_VERIFY_CODE, constants.STATE_SUBMIT:
		res, err = auth.forgotPasswordSubmit(ctx, body)
	default:
		err = variables.NewError(constants.ERR_REQUEST_INVALID)
	}
	return
}

func generateToken(id string) string {
	code := util.GenerateRandomString(36, nil)
	idEnc := base64.StdEncoding.EncodeToString([]byte(id))
	token := fmt.Sprintf("fins-%s-%s", code, idEnc)
	return token
}

func (s *Auth) forgotPasswordRequest(ctx context.Context, body model.ForgotPassword) (res *model.ForgotPasswordResponse, err error) {
	res = &model.ForgotPasswordResponse{
		IsValid: true,
	}

	user, err := repository.IBKUserRepo.GetByQuery(ctx, sql_builder.EqualQuery("email", body.Email))
	if err != nil {
		log.Error(err)
		return
	} else if user == nil {
		err = variables.NewError(constants.ERR_EMAIL_NOTFOUND)
		return
	} else if len(user.Email) < 1 {
		err = variables.NewError(constants.ERR_EMAIL_NOTFOUND)
		return
	}

	current := time.Now()

	request, err := repository.IBKUserRequestRepo.GetByQuery(ctx, sql_builder.EqualQuery("user_id", user.Id))
	if err != nil {
		log.Error(err)
		return
	} else if request != nil {
		isRemainTime := request.CreatedAt.Add(1 * time.Minute).After(current)
		if isRemainTime {
			res = &model.ForgotPasswordResponse{
				IsValid:         false,
				Code:            "ERR_REQUEST_IS_EXISTED",
				NextRequestTime: request.CreatedAt.Add(1 * time.Minute),
				// RemainTime:      current.Sub(aaaRequest.CreatedAt.Add(1 * time.Minute)),
			}
			return
		}
	}

	token := generateToken(user.Id)
	code := util.GenerateRandomString(6, util.NUMBER_RUNES)
	request = &model.IBKUserRequest{
		Base:         model.InitBaseModel("", ""),
		UserId:       user.Id,
		Method:       "forgot",
		Ip:           "",
		UserAgent:    body.UserAgent,
		Code:         code,
		Token:        token,
		InvalidCount: 0,
		ExpiredAt:    current.Add(30 * time.Minute),
		IsActive:     true,
	}

	if err = repository.IBKUserRequestRepo.Insert(ctx, *request); err != nil {
		log.Error(err)
		return
	}

	link := fmt.Sprintf("%s/forgot?token=%s", s.UIBaseUrl, token)
	log.Infof("link: %s", link)
	var emailBody string
	emailBody, err = mail.ParseTemplate("public/mail/forgot.html", struct {
		Link     string `json:"link"`
		Fullname string `json:"fullname"`
		Code     string `json:"code"`
	}{
		Link:     link,
		Fullname: user.Fullname,
		Code:     code,
	})
	if err != nil {
		log.Error(err)
		return
	}
	go func(email string, message string) {
		m := gomail.NewMessage()
		m.SetBody("text/html", message)
		if err = mail.Client.Send("[FinS] Forgot Password", email, m); err != nil {
			log.Error(err)
		}
	}(user.Email, emailBody)
	return
}

// verify token
func (s *Auth) forgotPasswordVerifyToken(ctx context.Context, body model.ForgotPassword) (res *model.ForgotPasswordResponse, err error) {
	request, err := repository.IBKUserRequestRepo.GetByQuery(ctx, sql_builder.EqualQuery("token", body.Token))
	if err != nil {
		log.Error(err)
		return
	} else if request == nil {
		err = variables.NewError(constants.ERR_REQUEST_INVALID)
		return
	}

	if request.ExpiredAt.Before(time.Now()) {
		err = variables.NewError(constants.ERR_REQUEST_IS_EXPIRED)
		return
	}

	res = &model.ForgotPasswordResponse{
		IsValid: true,
	}

	return
}

func (s *Auth) forgotPasswordSubmit(ctx context.Context, body model.ForgotPassword) (res *model.ForgotPasswordResponse, err error) {
	conditions := []sql_builder.QueryCondition{
		sql_builder.EqualQuery("token", body.Token),
		sql_builder.EqualQuery("is_active", true),
	}
	request, err := repository.IBKUserRequestRepo.GetByQuery(ctx, conditions...)
	if err != nil {
		log.Error(err)
		return
	} else if request == nil {
		err = variables.NewError(constants.ERR_REQUEST_INVALID)
		return
	}
	res = &model.ForgotPasswordResponse{
		IsValid: true,
	}
	if request.ExpiredAt.Before(time.Now()) {
		err = variables.NewError(constants.ERR_REQUEST_IS_EXPIRED)
		return
	} else if body.State == constants.STATE_VERIFY_TOKEN {
		return
	} else if body.Code != request.Code {
		request.InvalidCount = request.InvalidCount + 1
		if err = repository.IBKUserRequestRepo.Update(ctx, *request); err != nil {
			log.Error(err)
		}
		err = variables.NewError(constants.ERR_REQUEST_CODE_IS_INVALID)
		return
	} else if body.State == constants.STATE_VERIFY_CODE {
		return
	}
	salt := util.GenerateRandomString(8, nil)
	password := fmt.Sprintf("%s$%s", salt, hashSalt(salt, body.Password))
	user, err := repository.IBKUserRepo.GetById(ctx, request.UserId)
	if err != nil {
		log.Error(err)
		return
	} else if user == nil {
		err = variables.NewError(constants.ERR_REQUEST_INVALID)
		return
	}
	user.Password = password
	if err = repository.IBKUserRepo.Update(ctx, *user); err != nil {
		log.Error(err)
		return
	}
	request.IsActive = false
	if err = repository.IBKUserRequestRepo.Update(ctx, *request); err != nil {
		log.Error(err)
	}
	return
}

func (s *Auth) Logout(ctx context.Context, token, fingerprint string) (err error) {
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

func (s *Auth) RefreshAuthData(ctx context.Context, authUser *model.AuthUser, token string) (err error) {
	goAuthUser, err := goauth.GoAuthClient.Find(ctx, token)
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	} else if goAuthUser == nil {
		return errors.New("token is invalid")
	}
	if err = s.refreshAuthData(ctx, goAuthUser); err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	return
}

func (s *Auth) refreshAuthData(ctx context.Context, goAuthUser *goauth.AuthUser) (err error) {
	user, err := repository.IBKUserRepo.GetById(ctx, goAuthUser.UserId, sql_builder.QueryCondition{})
	if err != nil {
		log.ErrorContext(ctx, err)
		return err
	} else if user == nil {
		return errors.New("token is invalid")
	}
	tokenData := model.AuthUserData{
		BusinessUnitId: user.BusinessUnitId,
		Username:       user.Username,
		UserId:         user.Id,
		Level:          user.Level,
		Scopes:         user.Scopes,
	}
	goAuthUser.Data = tokenData
	if err := goauth.GoAuthClient.AddUserAuth(ctx, *goAuthUser); err != nil {
		log.ErrorContext(ctx, err)
		return err
	}
	return
}
