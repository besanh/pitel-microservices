package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
)

func RequestOttLabel(ctx context.Context, requestType, suffixUrl string, request model.ChatExternalLabelRequest) (result model.ChatExternalLabelResponse, err error) {
	chatLabelAfterRequest := model.ChatExternalLabelConvertedRequest{}
	if err = util.ParseAnyToAny(request, &chatLabelAfterRequest); err != nil {
		log.Error(err)
		return
	}
	chatLabelAfterRequest.UserId = request.ExternalUserId
	url := OTT_URL + "/ott/" + OTT_VERSION + "/crm/" + requestType + "/" + suffixUrl
	client := resty.New()
	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(url)

	if err != nil {
		log.Error(err)
		return
	}

	if err = json.Unmarshal([]byte(res.Body()), &result); err != nil {
		log.Error(err)
		return
	}

	if res.StatusCode() != 200 {
		err = errors.New(result.Message)
	}

	return
}
