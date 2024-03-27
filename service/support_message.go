package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
)

func (s *Message) sendMessageToOTT(ott model.SendMessageToOtt, attachment []*model.OttAttachments) (model.OttResponse, error) {
	var result model.OttResponse
	var body any
	var resMix model.SendMessageToOttWithAttachment
	resMix.Type = ott.Type
	resMix.EventName = ott.EventName
	resMix.AppId = ott.AppId
	resMix.OaId = ott.OaId
	resMix.UserIdByApp = ott.UserIdByApp
	resMix.Uid = ott.Uid
	resMix.SupporterId = ott.SupporterId
	resMix.SupporterName = ott.SupporterName
	resMix.Text = ott.Text
	resMix.Timestamp = ott.Timestamp
	resMix.MsgId = ott.MsgId

	if attachment != nil {
		resMix.Attachments = attachment
	}

	if err := util.ParseAnyToAny(resMix, &body); err != nil {
		return result, err
	}

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

func SendEventToManage(ctx context.Context, authUser *model.AuthUser, message model.Message, queueId string) (err error) {
	manageQueueUser, err := GetManageQueueUser(ctx, queueId)
	if err != nil {
		log.Error(err)
		return err
	} else if len(manageQueueUser.Id) < 1 {
		log.Error("queue " + queueId + " not found")
		err = errors.New("queue " + queueId + " not found")
		return err
	}

	// TODO: publish message to manager
	for s := range WsSubscribers.Subscribers {
		if s.Id == manageQueueUser.ManageId && s.Id != authUser.UserId {
			event := model.Event{
				EventName: variables.EVENT_CHAT[3],
				EventData: &model.EventData{
					Message: message,
				},
			}
			if err := PublishMessageToOne(manageQueueUser.ManageId, event); err != nil {
				log.Error(err)
				return err
			}
			break
		}
	}

	// TODO: publish to admin
	if ENABLE_PUBLISH_ADMIN {
		userUuids := []string{}
		for s := range WsSubscribers.Subscribers {
			if s.TenantId == message.TenantId && s.Level == "admin" && s.Id != authUser.UserId {
				userUuids = append(userUuids, s.Id)
			}
		}
		event := model.Event{
			EventName: variables.EVENT_CHAT[3],
			EventData: &model.EventData{
				Message: message,
			},
		}
		if err := PublishMessageToMany(userUuids, event); err != nil {
			log.Error(err)
			return err
		}
	}
	return
}
