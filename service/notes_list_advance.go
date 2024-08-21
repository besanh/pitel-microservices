package service

import (
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
)

/*
 * send event note to manager and admin subscribers
 */
func (s *NotesList) publishNotesListEventToManagerAndAdmin(authUser *model.AuthUser, manageQueueUser *model.ChatManageQueueUser, eventName string, note *model.NotesList) {
	var subscribers []*Subscriber
	var subscriberAdmins []string
	var subscriberManagers []string
	for sub := range WsSubscribers.Subscribers {
		if sub.TenantId == authUser.TenantId {
			subscribers = append(subscribers, sub)
			if sub.Level == "admin" {
				subscriberAdmins = append(subscriberAdmins, sub.Id)
			}
			if sub.Level == "manager" {
				subscriberManagers = append(subscriberManagers, sub.Id)
			}
		}
	}

	if manageQueueUser != nil {
		// Event to manager
		isExist := BinarySearchSlice(manageQueueUser.UserId, subscriberManagers)
		if isExist && len(manageQueueUser.UserId) > 0 {
			go PublishNotesListToOneUser(variables.EVENT_CHAT[eventName], manageQueueUser.UserId, subscribers, true, note)
		}
	}

	// Event to admin
	if ENABLE_PUBLISH_ADMIN && len(subscriberAdmins) > 0 {
		go PublishNotesListToManyUser(variables.EVENT_CHAT[eventName], subscriberAdmins, true, note)
	}
	go PublishNotesListToOneUser(variables.EVENT_CHAT[eventName], authUser.UserId, subscribers, true, note)
}
