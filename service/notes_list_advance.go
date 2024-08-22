package service

import (
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
)

/*
 * send event note to high level subscribers
 */
func (s *NotesList) publishNotesListEventToHighLevel(authUser *model.AuthUser, manageQueueUser *model.ChatManageQueueUser, eventName, eventDataType string, note *model.NotesList, levels ...string) {
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
			go PublishWsEventToOneUser(variables.EVENT_CHAT[eventName], eventDataType, manageQueueUser.UserId, subscribers["default"], true, note)
		}
	}

	// Event to admin
	if ENABLE_PUBLISH_ADMIN && len(subscribers[variables.ADMIN_LEVEL]) > 0 {
		go PublishWsEventToManyUser(variables.EVENT_CHAT[eventName], eventDataType, subscribers[variables.ADMIN_LEVEL], true, note)
	}
	go PublishWsEventToOneUser(variables.EVENT_CHAT[eventName], eventDataType, authUser.UserId, subscribers["default"], true, note)
}
