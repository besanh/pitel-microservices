package service

import (
	"context"
	"errors"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	ISubscriber interface {
		AddSubscriber(ctx context.Context, authUser *model.AuthUser, subscriber *Subscriber) (err error)
		GetSubscriber(id string) (*Subscriber, error)
		GetSubscribers() []Subscriber
		PublishMessageToSubscriber(ctx context.Context, id string, message any) error
	}
	SubscriberService struct{}
)

const (
	BSS_SUBSCRIBERS = "bss_subscribers"
)

var SubscriberServiceGlobal ISubscriber

func NewSubscriberService() *SubscriberService {
	if err := cache.RCache.Del([]string{BSS_SUBSCRIBERS}); err != nil {
		log.Error(err)
	}
	return &SubscriberService{}
}

func (s *SubscriberService) AddSubscriber(ctx context.Context, authUser *model.AuthUser, subscriber *Subscriber) (err error) {
	subscriber.Id = authUser.UserId
	subscriber.BusinessUnitId = authUser.BusinessUnitId
	subscriber.TenantId = authUser.TenantId
	subscriber.UserId = authUser.UserId
	subscriber.Username = authUser.Username
	subscriber.Level = authUser.Level
	subscriber.Services = authUser.Services
	subscriber.Source = authUser.Source
	subscriber.SubscribeAt = time.Now()

	WsSubscribers.AddSubscriber(ctx, subscriber)
	go func() {
		time.Sleep(1500 * time.Millisecond)
		message := map[string]any{
			"message": "init",
		}
		s.PublishMessageToSubscriber(ctx, subscriber.Id, message)
	}()

	return nil
}

func (s *SubscriberService) PublishMessageToSubscriber(ctx context.Context, id string, message any) error {
	err := PublishMessageToOne(id, message)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (s *SubscriberService) GetSubscriber(id string) (*Subscriber, error) {
	for s := range WsSubscribers.Subscribers {
		if s.Id == id {
			return s, nil
		}
	}
	return nil, errors.New("subscriber not found")
}

func (s *SubscriberService) GetSubscribers() []Subscriber {
	subscribers := []Subscriber{}
	for s := range WsSubscribers.Subscribers {
		subscribers = append(subscribers, *s)
	}
	return subscribers
}
