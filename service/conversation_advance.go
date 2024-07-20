package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"slices"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

/**
* Create in internal system, then create in external
 */
func PutLabelToConversation(ctx context.Context, authUser *model.AuthUser, labelType string, request model.ConversationLabelRequest) (labelId string, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}
	filter := model.ChatLabelFilter{
		TenantId:        authUser.TenantId,
		AppId:           request.AppId,
		OaId:            request.OaId,
		LabelType:       labelType,
		LabelName:       request.LabelName,
		IsSearchExactly: sql.NullBool{Bool: true, Valid: true},
		ExternalLabelId: request.ExternalLabelId,
	}
	_, chatLabelExists, err := repository.ChatLabelRepo.GetChatLabels(ctx, dbCon, filter, -1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*chatLabelExists) > 1 {
		log.Error("chat label " + request.LabelName + " already exists")
		err = errors.New("chat label " + request.LabelName + " already exists")
		return
	}
	if request.Action == "update" {
		if len(*chatLabelExists) == 0 {
			log.Error("chat label " + request.LabelName + " not found")
			err = errors.New("chat label " + request.LabelName + " not found")
			return
		}
	}

	// TODO: validate appId and oaId
	filterConnection := model.ChatConnectionAppFilter{
		TenantId:       authUser.TenantId,
		AppId:          request.AppId,
		OaId:           request.OaId,
		ConnectionType: labelType,
		Status:         "active",
	}
	_, connections, err := repository.ChatConnectionAppRepo.GetChatConnectionApp(ctx, dbCon, filterConnection, 1, 0)
	if err != nil {
		log.Error(err)
		return
	} else if len(*connections) == 0 {
		log.Error("connection with app_id: " + request.AppId + ", oa_id: " + request.OaId + " not found")
		err = errors.New("connection with app_id: " + request.AppId + ", oa_id: " + request.OaId + " not found")
		return
	}

	chatLabel := model.ChatLabel{
		Base:        model.InitBase(),
		TenantId:    authUser.TenantId,
		AppId:       request.AppId,
		OaId:        request.OaId,
		LabelName:   request.LabelName,
		LabelType:   labelType,
		LabelColor:  request.LabelColor,
		LabelStatus: true,
		CreatedBy:   authUser.UserId,
	}

	var externalLabelId string

	// TODO: create label zalo
	if labelType == "zalo" {
		if err = handleLabelZalo(ctx, labelType, request); err != nil {
			return
		}
		labelId = chatLabel.GetId()
	} else if labelType == "facebook" {
		externalLabelId, err = handleLabelFacebook(ctx, dbCon, labelType, chatLabel, request)
		if err != nil {
			return
		}
		// TODO: update label
		chatLabel.ExternalLabelId = externalLabelId
	}

	if request.Action == "create" {
		if err = repository.ChatLabelRepo.Insert(ctx, dbCon, chatLabel); err != nil {
			log.Error(err)
			return
		}
		if labelType == "zalo" {
			externalLabelId = chatLabel.GetId()
		}
	} else if request.Action == "update" {
		if labelType == "facebook" {
			// Get id to update
			filterTmp := model.ChatLabelFilter{
				TenantId:  authUser.TenantId,
				AppId:     request.AppId,
				OaId:      request.OaId,
				LabelName: request.LabelName,
				LabelType: labelType,
			}
			_, labelExist, errTmp := repository.ChatLabelRepo.GetChatLabels(ctx, dbCon, filterTmp, 1, 0)
			if errTmp != nil {
				log.Error(errTmp)
				return
			}
			if len(*labelExist) == 0 {
				log.Error("chat label " + request.LabelName + " not found")
				err = errors.New("chat label " + request.LabelName + " not found")
				return
			}
			(*labelExist)[0].ExternalLabelId = externalLabelId
			if err = repository.ChatLabelRepo.Update(ctx, dbCon, (*labelExist)[0]); err != nil {
				log.Error(err)
				return
			}
		} else {
			filterTmp := model.ChatLabelFilter{
				TenantId:  authUser.TenantId,
				AppId:     request.AppId,
				OaId:      request.OaId,
				LabelName: request.LabelName,
				LabelType: labelType,
			}
			_, labelExist, errTmp := repository.ChatLabelRepo.GetChatLabels(ctx, dbCon, filterTmp, 1, 0)
			if errTmp != nil {
				log.Error(errTmp)
				return
			}
			if len(*labelExist) == 0 {
				log.Error("chat label " + request.LabelName + " not found")
				err = errors.New("chat label " + request.LabelName + " not found")
				return
			}
			externalLabelId = (*labelExist)[0].Id
		}
	} else if request.Action == "delete" {
		if labelType == "facebook" {
			externalLabelId = request.ExternalLabelId
		} else {
			filterTmp := model.ChatLabelFilter{
				TenantId:  authUser.TenantId,
				AppId:     request.AppId,
				OaId:      request.OaId,
				LabelName: request.LabelName,
				LabelType: labelType,
			}
			_, labelExist, errTmp := repository.ChatLabelRepo.GetChatLabels(ctx, dbCon, filterTmp, 1, 0)
			if errTmp != nil {
				log.Error(errTmp)
				return
			}
			if len(*labelExist) == 0 {
				log.Error("chat label " + request.LabelName + " not found")
				err = errors.New("chat label " + request.LabelName + " not found")
				return
			}
			externalLabelId = (*labelExist)[0].Id
		}
	}

	// TODO: update label for conversation => use queue
	conversation, err := putConversation(ctx, authUser, externalLabelId, labelType, request)
	if err != nil {
		if request.Action == "create" {
			if err = repository.ChatLabelRepo.Delete(ctx, dbCon, chatLabel.GetId()); err != nil {
				log.Error(err)
				return
			}
		}
		return
	}

	labelId = chatLabel.GetId()

	// TODO: add event
	conversationConverted := &model.ConversationView{}
	if err = util.ParseAnyToAny(conversation, conversationConverted); err != nil {
		log.Error(err)
		return
	}

	if err = GetLabelsInfo(ctx, conversationConverted); err != nil {
		return
	}

	var subscribers []*Subscriber
	var subscriberAdmins []string
	// var subscriberManagers []string
	for s := range WsSubscribers.Subscribers {
		if s.TenantId == authUser.TenantId {
			subscribers = append(subscribers, s)
			if s.Level == "admin" {
				subscriberAdmins = append(subscriberAdmins, s.Id)
			}
			// if s.Level == "manager" {
			// 	subscriberManagers = append(subscriberManagers, s.Id)
			// }
		}
	}

	// TODO: publish event to user normal
	filterUserAllocate := model.AllocateUserFilter{
		TenantId:       authUser.TenantId,
		AppId:          request.AppId,
		OaId:           request.OaId,
		ConversationId: request.ConversationId,
	}
	_, userAllocate, err := repository.AllocateUserRepo.GetAllocateUsers(ctx, repository.DBConn, filterUserAllocate, 1, 0)
	if err != nil {
		log.Error(err)
		return
	} else if len(*userAllocate) > 0 {
		// Publish to user normal
		if request.Action == "create" || request.Action == "update" {
			if len(subscribers) > 0 {
				go PublishConversationToOneUser(variables.EVENT_CHAT["conversation_add_labels"], (*userAllocate)[0].UserId, subscribers, true, conversationConverted)
			}
		} else if request.Action == "delete" {
			if len(subscribers) > 0 {
				go PublishConversationToOneUser(variables.EVENT_CHAT["conversation_remove_labels"], (*userAllocate)[0].UserId, subscribers, true, conversationConverted)
			}
		}

		// TODO: get user manager then publish
		manageQueueUserFilter := model.ChatManageQueueUserFilter{
			TenantId: authUser.TenantId,
			QueueId:  (*userAllocate)[0].QueueId,
		}
		_, manageQueueUser, errTmp := repository.ManageQueueRepo.GetManageQueues(ctx, repository.DBConn, manageQueueUserFilter, 1, 0)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		} else if len(*manageQueueUser) > 0 {
			if request.Action == "create" || request.Action == "update" {
				if len(subscribers) > 0 {
					go PublishConversationToOneUser(variables.EVENT_CHAT["conversation_add_labels"], (*manageQueueUser)[0].UserId, subscribers, true, conversationConverted)
				}
			} else if request.Action == "delete" {
				if len(subscribers) > 0 {
					go PublishConversationToOneUser(variables.EVENT_CHAT["conversation_remove_labels"], (*manageQueueUser)[0].UserId, subscribers, true, conversationConverted)
				}
			}
		}
	}

	if ENABLE_PUBLISH_ADMIN {
		if request.Action == "create" || request.Action == "update" {
			if len(subscribers) > 0 {
				go PublishConversationToManyUser(variables.EVENT_CHAT["conversation_add_labels"], subscriberAdmins, true, conversationConverted)
			}
		} else if request.Action == "delete" {
			if len(subscribers) > 0 {
				go PublishConversationToManyUser(variables.EVENT_CHAT["conversation_remove_labels"], subscriberAdmins, true, conversationConverted)
			}
		}
	}

	return
}

func handleLabelZalo(ctx context.Context, labelType string, request model.ConversationLabelRequest) (err error) {
	zaloRequest := model.ChatExternalLabelRequest{
		AppId:          request.AppId,
		OaId:           request.OaId,
		ExternalUserId: request.ExternalUserId,
		TagName:        request.LabelName,
	}
	var externalUrl string
	if request.Action == "create" || request.Action == "update" {
		externalUrl = "create-label-customer"
	} else if request.Action == "update" {

	} else if request.Action == "delete" {
		externalUrl = "remove-label-customer"
	}

	// TODO: because zalo not return id so we don't need to use it for updating label
	if slices.Contains([]string{"create", "update", "delete"}, request.Action) {
		_, errTmp := RequestOttLabel(ctx, labelType, externalUrl, zaloRequest)
		if errTmp != nil {
			log.Error(errTmp)
			err = errTmp
			return
		}
	}

	return
}

func handleLabelFacebook(ctx context.Context, dbCon sqlclient.ISqlClientConn, labelType string, chatLabel model.ChatLabel, request model.ConversationLabelRequest) (externalLabelId string, err error) {
	if labelType == "facebook" {
		labelType = "face"
	}

	facebookRequest := model.ChatExternalLabelRequest{
		AppId:          request.AppId,
		OaId:           request.OaId,
		ExternalUserId: request.ExternalUserId,
		LabelId:        request.ExternalLabelId,
		TagName:        request.LabelName,
	}

	var externalUrl string

	if request.Action == "create" || request.Action == "update" {
		// TODO: if label does not exist, create new label, then associate to conversation)
		externalUrl = "create-label"
		externalCreateLabelResponse, errTmp := RequestOttLabel(ctx, labelType, externalUrl, facebookRequest)
		if errTmp != nil {
			log.Error(errTmp)
			if err = repository.ChatLabelRepo.Delete(ctx, dbCon, chatLabel.GetId()); err != nil {
				log.Error(err)
				return externalLabelId, err
			}
			return externalLabelId, errTmp
		}
		if len(externalCreateLabelResponse.Id) > 0 {
			externalLabelId = externalCreateLabelResponse.Id
		} else {
			err = errors.New("external label id not found")
			log.Error("external label id not found")
			return
		}

		// TODO: associating a label to oa
		facebookAssociateRequest := model.ChatExternalLabelRequest{
			AppId:          request.AppId,
			OaId:           request.OaId,
			LabelId:        externalLabelId,
			ExternalUserId: request.ExternalUserId,
			TagName:        request.LabelName,
		}
		externalUrl = "associate-label"
		_, errTmp = RequestOttLabel(ctx, labelType, externalUrl, facebookAssociateRequest)
		if errTmp != nil {
			log.Error(errTmp)
			if err = repository.ChatLabelRepo.Delete(ctx, dbCon, chatLabel.GetId()); err != nil {
				log.Error(err)
				return externalLabelId, err
			}
			err = errTmp
			return externalLabelId, errTmp
		}
	} else if request.Action == "delete" {
		externalUrl = "remove-label"
		_, errTmp := RequestOttLabel(ctx, labelType, externalUrl, facebookRequest)
		if errTmp != nil {
			log.Error(errTmp)
			if err = repository.ChatLabelRepo.Delete(ctx, dbCon, chatLabel.GetId()); err != nil {
				log.Error(err)
				return externalLabelId, err
			}
			err = errTmp
			return externalLabelId, errTmp
		}
	} else {
		err = errors.New("invalid action")
		log.Error("invalid action")
	}
	return
}

func putConversation(ctx context.Context, authUser *model.AuthUser, labelId, labelType string, request model.ConversationLabelRequest) (conversationExist *model.Conversation, err error) {
	conversationExist, err = repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, request.AppId, request.ConversationId)
	if err != nil {
		log.Error(err)
		return
	} else if len(conversationExist.ConversationId) < 1 {
		log.Errorf("conversation %s not found", request.ConversationId)
		return
	}

	result, err := UpdateConversationLabelList(conversationExist.Label, labelType, request.Action, labelId)
	if err != nil {
		return
	}
	conversationExist.Label = result

	if err = PublishPutConversationToChatQueue(ctx, *conversationExist); err != nil {
		log.Error(err)
		return
	}
	return
}

func UpdateConversationLabelList(existLabels json.RawMessage, labelType string, action string, labelId string) (result []byte, err error) {
	objmap := []any{}
	labelsExist := []any{}
	if err = json.Unmarshal([]byte(existLabels), &labelsExist); err != nil {
		log.Error(err)
		return
	}

	// TODO: because zalo only assign one label for one conversation
	if labelType == "facebook" {
		for _, item := range labelsExist {
			tmp := map[string]string{}
			if err = util.ParseAnyToAny(item, &tmp); err != nil {
				log.Error(err)
				continue
			}
			if action == "delete" && tmp["label_id"] == labelId {
				continue
			}
			if len(tmp["label_id"]) > 0 {
				isExist := checkItemExist(objmap, tmp)
				if !isExist {
					objmap = append(objmap, map[string]any{
						"label_id": tmp["label_id"],
					})
				}
			}
		}
	}

	if action == "create" || action == "update" {
		if len(labelId) > 0 {
			if len(labelsExist) > 0 {
				for _, item := range labelsExist {
					id := item.(map[string]any)["label_id"].(string)
					if id == labelId {
						continue
					}
					isExist := checkItemExist(objmap, map[string]string{"label_id": id})
					if !isExist {
						objmap = append(objmap, map[string]any{
							"label_id": id,
						})
					}
				}
				isExist := checkItemExist(objmap, map[string]string{"label_id": labelId})
				if !isExist {
					objmap = append(objmap, map[string]any{
						"label_id": labelId,
					})
				}
			} else {
				objmap = append(objmap, map[string]any{
					"label_id": labelId,
				})
			}
		}
	} else if action == "delete" {
		if labelType == "zalo" {
			for _, item := range labelsExist {
				tmp := map[string]string{}
				if err = util.ParseAnyToAny(item, &tmp); err != nil {
					log.Error(err)
					continue
				}
				if tmp["label_id"] == labelId {
					continue
				}
				isExist := checkItemExist(objmap, tmp)
				if !isExist {
					objmap = append(objmap, map[string]any{
						"label_id": tmp["label_id"],
					})
				}
			}
		}
	}
	log.Info("label_id: ", labelId)

	result, err = json.Marshal(objmap)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func checkItemExist(objmap []any, tmp map[string]string) (isExist bool) {
	for _, item := range objmap {
		if item.(map[string]any)["label_id"] == tmp["label_id"] {
			return true
		}
	}
	return
}

/*
 * send event conversation to manager and admin subscribers
 */
func (s *Conversation) publishConversationEventToManagerAndAdmin(authUser *model.AuthUser, manageQueueUser *model.ChatManageQueueUser, eventName string, conversationConverted *model.ConversationView) {
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
			go PublishConversationToOneUser(variables.EVENT_CHAT[eventName], manageQueueUser.UserId, subscribers, true, conversationConverted)
		}
	}

	// Event to admin
	if ENABLE_PUBLISH_ADMIN && len(subscriberAdmins) > 0 {
		go PublishConversationToManyUser(variables.EVENT_CHAT[eventName], subscriberAdmins, true, conversationConverted)
	}
	go PublishConversationToOneUser(variables.EVENT_CHAT[eventName], authUser.UserId, subscribers, true, conversationConverted)
}
