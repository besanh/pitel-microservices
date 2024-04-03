package common

import (
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
)

func PostOttAccount(ottDomain, ottVersion string, chatApp model.ChatApp, connectionApp model.ChatConnectionApp) error {
	accountInfo := model.OttAccount{}
	if chatApp.InfoApp.Facebook.Status == "active" {
		accountInfo.Type = "face"
		accountInfo.AppId = chatApp.InfoApp.Facebook.AppId
		accountInfo.AppName = chatApp.InfoApp.Facebook.AppName
		accountInfo.AppSecret = chatApp.InfoApp.Facebook.AppSecret
		accountInfo.OaId = connectionApp.OaInfo.Facebook[0].OaId
		accountInfo.OaName = connectionApp.OaInfo.Facebook[0].OaName
		accountInfo.AccessToken = connectionApp.OaInfo.Facebook[0].AccessToken
		accountInfo.Status = chatApp.InfoApp.Facebook.Status
	} else if chatApp.InfoApp.Zalo.Status == "active" {
		accountInfo.Type = "zalo"
		accountInfo.AppId = chatApp.InfoApp.Zalo.AppId
		accountInfo.AppName = chatApp.InfoApp.Zalo.AppName
		accountInfo.AppSecret = chatApp.InfoApp.Zalo.AppSecret
		accountInfo.OaId = connectionApp.OaInfo.Zalo[0].OaId
		accountInfo.OaName = connectionApp.OaInfo.Zalo[0].OaName
		accountInfo.Status = chatApp.InfoApp.Zalo.Status
	}

	body := map[string]string{
		"type":         accountInfo.Type,
		"app_id":       accountInfo.AppId,
		"app_name":     accountInfo.AppName,
		"app_secret":   accountInfo.AppSecret,
		"oa_id":        accountInfo.OaId,
		"oa_name":      accountInfo.OaName,
		"status":       accountInfo.Status,
		"access_token": accountInfo.AccessToken,
	}

	log.Info("post ott account: ", body)

	url := ottDomain + "/ott/" + ottVersion + "/account"
	client := resty.New()

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		Post(url)
	if err != nil {
		return err
	}

	if res.StatusCode() == 200 {
		return nil
	} else {
		return errors.New("create app error")
	}
}
