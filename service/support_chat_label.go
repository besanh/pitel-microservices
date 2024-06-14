package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
)

func (s *ChatLabel) RequestZaloLabel(ctx context.Context, suffixUrl string, request model.ChatExtenalLabelRequest) (err error) {
	url := OTT_URL + "/ott/" + OTT_VERSION + "/crm/zalo/" + suffixUrl
	client := resty.New()

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(url)
	if err != nil {
		log.Error(err)
		return err
	}

	var result model.ChatZaloLabelResponse

	if err := json.Unmarshal([]byte(res.Body()), &result); err != nil {
		log.Error(err)
		return err
	}
	if res.StatusCode() == 200 {
		return nil
	} else {
		err = errors.New(result.Message)
		return err
	}
}

func (s *ChatLabel) RequestFacebookLabel(ctx context.Context, suffixUrl string, request model.ChatExtenalLabelRequest) (err error) {
	url := OTT_URL + "/ott/" + OTT_VERSION + "/crm/face/" + suffixUrl
	client := resty.New()

	var res *resty.Response
	if suffixUrl != "" {
		res, err = client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(request).
			Post(url)
	} else {
		res, err = client.R().
			SetHeader("Content-Type", "application/json").
			Delete(url)
	}

	if err != nil {
		log.Error(err)
		return err
	}

	var result model.ChatZaloLabelResponse

	if err := json.Unmarshal([]byte(res.Body()), &result); err != nil {
		log.Error(err)
		return err
	}
	if res.StatusCode() == 200 {
		return nil
	} else {
		err = errors.New(result.Message)
		return err
	}
}
