package service

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
)

func (s *Message) sendMessageToOTT(ctx context.Context, ott model.SendMessageToOtt) (model.OttResponse, error) {
	var result model.OttResponse
	var body any
	if err := util.ParseAnyToAny(ott, &body); err != nil {
		return result, err
	}

	url := s.OttSendMessageUrl
	client := resty.New()

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		Post(url)
	if err != nil {
		return result, err
	}
	if err := util.ParseAnyToAny(res.Body(), &result); err != nil {
		return result, err
	}

	return result, nil
}
