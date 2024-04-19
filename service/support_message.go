package service

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

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

	url := OTT_URL + "/ott/" + OTT_VERSION + "/crm"
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
	var subscribers []*Subscriber
	var subscriberAdmins []string
	for s := range WsSubscribers.Subscribers {
		if s.TenantId == manageQueueUser.TenantId {
			subscribers = append(subscribers, s)
			if s.Level == "admin" {
				subscriberAdmins = append(subscriberAdmins, s.Id)
			}
		}
	}

	go PublishMessageToOneUser(variables.EVENT_CHAT["message_created"], manageQueueUser.ManageId, subscribers, &message)

	// TODO: publish to admin
	if ENABLE_PUBLISH_ADMIN {
		go PublishMessageToManyUser(variables.EVENT_CHAT["message_created"], subscriberAdmins, &message)
	}
	return
}

func PublishConversationToOneUser(eventType string, subscriber string, subscribers []*Subscriber, isNew bool, conversation *model.Conversation) {
	var wg sync.WaitGroup
	if isNew && conversation != nil && len(subscriber) > 0 {
		event := model.Event{
			EventName: eventType,
			EventData: &model.EventData{
				Conversation: conversation,
			},
		}

		isExist := BinarySearchSubscriber(subscriber, subscribers)
		if isExist {
			wg.Add(1)
			var mu sync.Mutex
			mu.Lock()
			go func(userUuid string, event model.Event) {
				defer wg.Done()
				if err := PublishMessageToOne(userUuid, event); err != nil {
					log.Error(err)
					return
				}
			}(subscriber, event)
			mu.Unlock()
		}
	}
	wg.Wait()
}

func PublishConversationToManyUser(eventType string, subscribers []string, isNew bool, conversation *model.Conversation) {
	var wg sync.WaitGroup
	if isNew && conversation != nil && len(subscribers) > 0 {
		event := model.Event{
			EventName: eventType,
			EventData: &model.EventData{
				Conversation: conversation,
			},
		}
		wg.Add(1)
		var mu sync.Mutex
		mu.Lock()
		go func(userUuids []string, event model.Event) {
			defer wg.Done()
			if err := PublishMessageToMany(userUuids, event); err != nil {
				log.Error(err)
				return
			}
		}(subscribers, event)
		mu.Unlock()
	}
	wg.Wait()
}

func PublishMessageToOneUser(eventType string, subscriber string, subscribers []*Subscriber, message *model.Message) {
	var wg sync.WaitGroup
	if message != nil && len(subscriber) > 0 {
		event := model.Event{
			EventName: eventType,
			EventData: &model.EventData{
				Message: message,
			},
		}
		isExist := BinarySearchSubscriber(subscriber, subscribers)
		if isExist {
			wg.Add(1)
			var mu sync.Mutex
			mu.Lock()
			go func(userUuid string, event model.Event) {
				defer wg.Done()
				if err := PublishMessageToOne(userUuid, event); err != nil {
					log.Error(err)
					return
				}
			}(subscriber, event)
			mu.Unlock()
		}
	}
	wg.Wait()
}

func PublishMessageToManyUser(eventType string, subscribers []string, message *model.Message) {
	var wg sync.WaitGroup
	if len(subscribers) > 0 && message != nil {
		event := model.Event{
			EventName: eventType,
			EventData: &model.EventData{
				Message: message,
			},
		}

		wg.Add(1)
		var mu sync.Mutex
		mu.Lock()
		go func(userUuids []string, event model.Event) {
			defer wg.Done()
			if err := PublishMessageToMany(userUuids, event); err != nil {
				log.Error(err)
				return
			}
		}(subscribers, event)
		mu.Unlock()
	}
	wg.Wait()
}
