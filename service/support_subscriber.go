package service

import (
	"context"
	"encoding/json"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
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
	if err := cache.RCache.HSetRaw(ctx, SUBSCRIBERS_LIST_USER, sub.UserId, string(jsonByte)); err != nil {
		log.Error(err)
	}

}

func (s *Subscribers) DeleteSubscriber(sub *Subscriber) {
	s.SubscribersMu.Lock()
	defer s.SubscribersMu.Unlock()

	delete(s.Subscribers, sub)
	if err := cache.RCache.HDel(SUBSCRIBERS_LIST_USER+"_"+sub.UserId, sub.UserId); err != nil {
		log.Error(err)
	}
}
