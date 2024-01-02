package service

import (
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
)

var WsSubscribers *Subscribers

func (s *Subscribers) AddSubscriber(sub *Subscriber) {
	s.SubscribersMu.Lock()
	defer s.SubscribersMu.Unlock()

	s.Subscribers[sub] = struct{}{}
	var tmp []any
	if err := util.ParseAnyToAny(sub, &tmp); err != nil {
		log.Error(err)
	}
	if err := cache.RCache.HSet(SUBSCRIBERS_LIST_USER+"_"+sub.UserId, tmp); err != nil {
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
