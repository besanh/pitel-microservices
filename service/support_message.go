package service

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
)

func (s *Message) sendMessageToOTT(ott model.SendMessageToOtt, attachment *[]model.OttAttachments) (model.OttResponse, error) {
	var result model.OttResponse
	var body any
	if err := util.ParseAnyToAny(ott, &body); err != nil {
		return result, err
	}
	// if attachment != nil {
	// 	body = append(body.([]any), attachment)
	// }

	url := OTT_URL + "/ott/v1/crm"
	client := resty.New()

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		Post(url)
	if err != nil {
		return result, err
	}

	if err := json.Unmarshal([]byte(res.Body()), &result); err != nil {
		return result, err
	}
	if res.StatusCode() == 200 {
		return result, nil
	} else {
		return result, errors.New(result.Message)
	}
}
