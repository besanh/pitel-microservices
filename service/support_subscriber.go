package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"golang.org/x/exp/slices"
)

var WsSubscribers *Subscribers

func (s *Subscribers) AddSubscriber(ctx context.Context, sub *Subscriber) {
	s.SubscribersMu.Lock()
	defer s.SubscribersMu.Unlock()

	s.Subscribers[sub] = struct{}{}
	jsonByte, err := json.Marshal(&sub)
	if err != nil {
		log.Error(err)
	}
	if err := cache.RCache.HSetRaw(ctx, BSS_SUBSCRIBERS, sub.UserId, string(jsonByte)); err != nil {
		log.Error(err)
	}

}

func (s *Subscribers) DeleteSubscriber(sub *Subscriber) {
	s.SubscribersMu.Lock()
	defer s.SubscribersMu.Unlock()

	delete(s.Subscribers, sub)
	if err := cache.RCache.HDel(BSS_SUBSCRIBERS, sub.UserId); err != nil {
		log.Error(err)
	}
}

func PublishMessageToOne(id string, message any) error {
	WsSubscribers.SubscribersMu.Lock()
	defer WsSubscribers.SubscribersMu.Unlock()
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	WsSubscribers.PublishLimiter.Wait(context.Background())
	isExisted := false
	for s := range WsSubscribers.Subscribers {
		if s.Id == id {
			isExisted = true
			select {
			case s.Message <- msgBytes:
			default:
				go s.CloseSlow()
			}
		}
	}
	if !isExisted {
		return errors.New("subscriber is not existed")
	}
	return nil
}

func PublishMessageToMany(ids []string, message any) error {
	WsSubscribers.SubscribersMu.Lock()
	defer WsSubscribers.SubscribersMu.Unlock()
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	WsSubscribers.PublishLimiter.Wait(context.Background())
	isExisted := false
	for s := range WsSubscribers.Subscribers {
		if slices.Contains(ids, s.Id) {
			isExisted = true
			select {
			case s.Message <- msgBytes:
			default:
				go s.CloseSlow()
			}
		}
	}
	if !isExisted {
		return errors.New("subscriber is not existed")
	}
	return nil
}
