package common

import (
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/model"
)

func PostOttAccount(ottDomain string, chatApp model.ChatApp) error {
	accountInfo := model.OttAccount{}
	if chatApp.InfoApp.Facebook.Status == "active" {
		accountInfo.Type = "facebook"
		accountInfo.AppId = chatApp.InfoApp.Facebook.AppId
		accountInfo.AppName = chatApp.InfoApp.Facebook.AppName
		accountInfo.AppSecret = chatApp.InfoApp.Facebook.AppSecret
		accountInfo.OaId = chatApp.InfoApp.Facebook.OaId
		accountInfo.OaName = chatApp.InfoApp.Facebook.OaName
		accountInfo.Status = chatApp.InfoApp.Facebook.Status
	} else if chatApp.InfoApp.Zalo.Status == "active" {
		accountInfo.Type = "zalo"
		accountInfo.AppId = chatApp.InfoApp.Zalo.AppId
		accountInfo.AppName = chatApp.InfoApp.Zalo.AppName
		accountInfo.AppSecret = chatApp.InfoApp.Zalo.AppSecret
		accountInfo.OaId = chatApp.InfoApp.Zalo.OaId
		accountInfo.OaName = chatApp.InfoApp.Zalo.OaName
		accountInfo.Status = chatApp.InfoApp.Zalo.Status
	}

	body := map[string]string{
		"type":       accountInfo.Type,
		"app_id":     accountInfo.AppId,
		"app_name":   accountInfo.AppName,
		"app_secret": accountInfo.AppSecret,
		"oa_id":      accountInfo.OaId,
		"oa_name":    accountInfo.OaName,
		"status":     accountInfo.Status,
	}

	url := ottDomain + "/ott/v1/account"
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

// func UpdateOttAccount(ottDomain string)
