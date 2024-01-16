package service

import (
	"context"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
)

func (s *Message) sendMessageToOTT(ctx context.Context, ott model.SendMessageToOtt) (int, model.OttResponse, error) {
	var result model.OttResponse
	var body any
	if err := util.ParseAnyToAny(ott, &body); err != nil {
		return http.StatusBadRequest, result, err
	}

	url := s.OttSendMessageUrl
	client := resty.New()

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		SetResult(&result).
		Post(url)
	if err != nil {
		return res.StatusCode(), result, err
	}

	return res.StatusCode(), result, nil
}
