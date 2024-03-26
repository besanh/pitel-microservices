package common

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/model"
)

func CheckUserExist(id string, Users []string) bool {
	for _, User := range Users {
		if User == id {
			return true
		}
	}
	return false
}

func GetUserAuthenticated(crmUrl, token, userId string) (*model.AuthUserInfo, error) {
	var result model.AuthUserInfo
	url := crmUrl + "/v1/crm/user-crm/" + userId
	client := resty.New()

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		Get(url)
	if err != nil {
		return &result, err
	}

	if err := json.Unmarshal([]byte(res.Body()), &result); err != nil {
		return &result, err
	}

	return &result, nil
}
