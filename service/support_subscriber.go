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
		return errors.New("publish to one -> subscriber " + id + " is not existed")
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
		return errors.New("publish to many -> subscriber is not existed")
	}
	return nil
}

func BinarySearchSubscriber(userId string, subscribers []*Subscriber) (isExist bool) {
	for s := range subscribers {
		if subscribers[s].Id == userId {
			return true
		}
	}

	return
}

func BinarySearchSlice(userId string, subscribers []string) (isExist bool) {
	low := 0
	high := len(subscribers) - 1
	mid := -1
	for low <= high {
		mid = (low + high) / 2
		if subscribers[mid] == userId {
			return true
		} else {
			high = mid - 1
		}
	}
	if mid != -1 {
		isExist = true
	}
	return
}

// Search with compare
// func bSearch(people []Person, target Person, low int, high int) int {
// 	if low > high {
// 		return -1
// 	}

// 	mid := (low + high) / 2
// 	if people[mid] == target {
// 		return mid
// 	} else if people[mid].Age < target.Age {
// 		return bSearch(people, target, mid+1, high)
// 	} else {
// 		return bSearch(people, target, low, mid-1)
// 	}
// }
