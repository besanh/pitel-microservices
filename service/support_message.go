package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
)

func sendMessageToOTT(ott model.SendMessageToOtt, attachment []*model.OttAttachments) (result model.OttResponse, err error) {
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

	if err = util.ParseAnyToAny(resMix, &body); err != nil {
		return
	}

	url := OTT_URL + "/ott/" + OTT_VERSION + "/crm"
	client := resty.New().
		SetTimeout(2 * time.Minute)

	res, err := client.R().
		SetHeader("Content-Type", "application/json").
		// SetHeader("Authorization", "Bearer "+token).
		SetBody(body).
		Post(url)
	if err != nil {
		return
	}

	if err = json.Unmarshal([]byte(res.Body()), &result); err != nil {
		return
	}
	log.Info("result: ", result)
	if res.StatusCode() != 200 {
		err = errors.New(result.Message)
	}
	return
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

	go PublishMessageToOneUser(variables.EVENT_CHAT["message_created"], manageQueueUser.UserId, subscribers, &message)

	// TODO: publish to admin
	if ENABLE_PUBLISH_ADMIN {
		go PublishMessageToManyUser(variables.EVENT_CHAT["message_created"], subscriberAdmins, &message)
	}
	return
}

func PublishConversationToOneUser(eventType string, subscriber string, subscribers []*Subscriber, isNew bool, conversation *model.ConversationView) {
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

func PublishConversationToManyUser(eventType string, subscribers []string, isNew bool, conversation *model.ConversationView) {
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

func PublishWsEventToOneUser(eventType, eventDataType string, subscriber string, subscribers []string, isNew bool, data any) {
	var wg sync.WaitGroup
	if isNew && data != nil && len(subscriber) > 0 {
		event := model.WsEvent{
			EventName: eventType,
			EventData: map[string]any{
				eventDataType: data,
			},
		}

		isExist := BinarySearchSlice(subscriber, subscribers)
		if isExist {
			wg.Add(1)
			var mu sync.Mutex
			mu.Lock()
			go func(userUuid string, event model.WsEvent) {
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

func PublishWsEventToManyUser(eventType, eventDataType string, subscribers []string, isNew bool, data any) {
	var wg sync.WaitGroup
	if isNew && data != nil && len(subscribers) > 0 {
		event := model.WsEvent{
			EventName: eventType,
			EventData: map[string]any{
				eventDataType: data,
			},
		}
		wg.Add(1)
		var mu sync.Mutex
		mu.Lock()
		go func(userUuids []string, event model.WsEvent) {
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

/*
 * send event note to high level subscribers
 */
func PublishEventToHighLevel(authUser *model.AuthUser, manageQueueUser *model.ChatManageQueueUser, eventName, eventDataType string, data any, levels ...string) {
	// default send it to admin
	if len(levels) < 1 {
		levels = append(levels, variables.ADMIN_LEVEL)
	}

	convertedData, err := renameFields(data)
	if err != nil {
		log.Error(err)
	}

	subscribers := make(map[string][]string)
	for sub := range WsSubscribers.Subscribers {
		if sub.TenantId == authUser.TenantId {
			if _, ok := subscribers["default"]; !ok {
				subscribers["default"] = []string{}
			}
			subscribers["default"] = append(subscribers["default"], sub.Id)
			// high level subscriber
			for _, level := range levels {
				if _, ok := subscribers[level]; !ok {
					subscribers[level] = []string{}
				}
				if sub.Level == level {
					subscribers[level] = append(subscribers[level], sub.Id)
				}
			}
		}
	}

	if manageQueueUser != nil {
		// Event to manager
		isExist := BinarySearchSlice(manageQueueUser.UserId, subscribers[variables.MANAGER_LEVEL])
		if isExist && len(manageQueueUser.UserId) > 0 {
			go PublishWsEventToOneUser(variables.EVENT_CHAT[eventName], eventDataType, manageQueueUser.UserId, subscribers["default"], true, convertedData)
		}
	}

	// Event to admin
	if ENABLE_PUBLISH_ADMIN && len(subscribers[variables.ADMIN_LEVEL]) > 0 {
		go PublishWsEventToManyUser(variables.EVENT_CHAT[eventName], eventDataType, subscribers[variables.ADMIN_LEVEL], true, convertedData)
	}
	go PublishWsEventToOneUser(variables.EVENT_CHAT[eventName], eventDataType, authUser.UserId, subscribers["default"], true, convertedData)
}

func renameFields(data any) (any, error) {
	// Use reflection to inspect the struct
	val := reflect.ValueOf(data)
	kind := val.Kind()
	switch kind {
	case reflect.Struct:
	case reflect.Map:
	default:
		return data, nil
	}
	// Convert struct to a map
	dataMap := make(map[string]any)
	if err := util.ParseAnyToAny(&data, &dataMap); err != nil {
		return data, err
	}
	if _, ok := dataMap["CreatedAt"]; ok {
		dataMap["created_at"] = dataMap["CreatedAt"]
	}
	if _, ok := dataMap["UpdatedAt"]; ok {
		dataMap["updated_at"] = dataMap["UpdatedAt"]
	}

	return dataMap, nil
}

func SendPushNotification(payload model.NotifyPayload) (err error) {
	if !ENABLE_NOTIFY_CHAT {
		return
	}
	res, err := resty.New().SetTimeout(10*time.Second).R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(API_PUSH_NOTIFICATION_SECRET).
		SetBody(payload).
		Post(API_PUSH_NOTIFICATION)
	if err != nil {
		return
	}
	if res.IsError() {
		err = fmt.Errorf("error %v, status code %d, err %v", res.Error(), res.StatusCode(), err)
		return
	}
	return
}
