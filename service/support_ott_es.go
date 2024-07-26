package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func (s *OttMessage) UpSertConversation(ctx context.Context, connectionId, conversationId string, data model.OttMessage) (conversation model.Conversation, isNew bool, err error) {
	externalConversationId := GenerateConversationId(data.AppId, data.OaId, data.ExternalUserId)
	conversation = model.Conversation{
		TenantId:               data.TenantId,
		ConversationId:         conversationId,
		ExternalConversationId: externalConversationId,
		AppId:                  data.AppId,
		ConversationType:       data.MessageType,
		Username:               data.Username,
		Avatar:                 data.Avatar,
		OaId:                   data.OaId,
		ShareInfo:              data.ShareInfo,
		ExternalUserId:         data.ExternalUserId,
		CreatedAt:              time.Now().Format(time.RFC3339),
	}
	shareInfo := data.ShareInfo
	isExisted := false
	conversationExist := model.ConversationView{}
	// TODO: improve by caching
	_, conversationExists, err := repository.ConversationESRepo.GetConversations(ctx, data.TenantId, ES_INDEX_CONVERSATION, model.ConversationFilter{
		TenantId:               conversation.TenantId,
		ExternalConversationId: []string{externalConversationId},
	}, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if len(*conversationExists) > 0 {
		conversationExist = (*conversationExists)[0]
		conversationId = conversationExist.ConversationId
		conversation.TenantId = conversationExist.TenantId
		conversation.ConversationId = conversationExist.ConversationId
		conversation.ExternalConversationId = conversationExist.ExternalConversationId
		conversation.ConversationType = conversationExist.ConversationType
		conversation.AppId = conversationExist.AppId
		conversation.OaId = conversationExist.OaId
		conversation.OaName = conversationExist.OaName
		conversation.OaAvatar = conversationExist.OaAvatar
		conversation.ExternalUserId = conversationExist.ExternalUserId
		conversation.Username = conversationExist.Username
		conversation.Avatar = conversationExist.Avatar
		conversation.Major = conversationExist.Major
		conversation.Following = conversationExist.Following
		conversation.Labels = conversationExist.Labels
		conversation.IsDone = false
		conversation.IsDoneBy = ""
		isDoneAt, _ := time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
		conversation.IsDoneAt = isDoneAt
		conversation.CreatedAt = conversationExist.CreatedAt
		conversation.UpdatedAt = time.Now().Format(time.RFC3339)

		conversation.ShareInfo = shareInfo
		// if len(connectionId) > 0 {
		// 	conversation, err = CacheConnection(ctx, connectionId, conversation)
		// 	if err != nil {
		// 		log.Error(err)
		// 		return conversation, isNew, err
		// 	}
		// }

		// TODO: update conversation => use queue consumer
		if !data.IsEcho {
			conversationQueue := model.ConversationQueue{
				DocId:        conversationId,
				Conversation: conversation,
			}
			if err = PublishPutConversationToChatQueue(ctx, conversationQueue); err != nil {
				log.Error(err)
				return
			}
		}

		isExisted = true
		return
	} else {
		conversation.Labels = []byte("[]")
	}

	if !isExisted {
		err = InsertConversation(ctx, conversationId, connectionId, conversation)
		if err != nil {
			// log.Error(err) // remove for duplicating
			return
		}
		if len(connectionId) > 0 {
			conversation, err = CacheConnection(ctx, connectionId, conversation)
			if err != nil {
				log.Error(err)
				return
			}
		}
		if err = cache.RCache.Set(CONVERSATION+"_"+conversationId, conversation, CONVERSATION_EXPIRE); err != nil {
			log.Error(err)
			return
		}
		isNew = true
	}

	return
}

func InsertConversation(ctx context.Context, docId, connectionId string, conversation model.Conversation) (err error) {
	if len(connectionId) > 0 {
		conversation, err = CacheConnection(ctx, connectionId, conversation)
		if err != nil {
			log.Error(err)
			return
		}
	}
	tmpBytes, err := json.Marshal(conversation)
	if err != nil {
		log.Error(err)
		return
	}

	esDoc := map[string]any{}
	if err = json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return
	}
	isExisted, errTmp := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX_CONVERSATION, conversation.TenantId)
	if errTmp != nil {
		log.Error(errTmp)
		err = errTmp
		return
	} else if !isExisted {
		if err = repository.ESRepo.CreateAlias(ctx, ES_INDEX_CONVERSATION, conversation.TenantId); err != nil {
			log.Error(err)
			return
		}
	}

	if err = repository.ESRepo.InsertLog(ctx, conversation.TenantId, ES_INDEX_CONVERSATION, conversation.AppId, docId, esDoc); err != nil {
		log.Error(err)
		return
	}

	return
}

/**
* Update ES and Cache
* API get conversation can get from redis, and here can caching to descrese the number of api calls to ES
 */
func (s *OttMessage) UpdateESAndCache(ctx context.Context, tenantId, appId, oaId, conversationId, connectionId string, shareInfo model.ShareInfo) (err error) {
	var isUpdate bool
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, tenantId, ES_INDEX_CONVERSATION, appId, conversationId)
	if err != nil {
		log.Error(err)
		return
	} else if len(conversationExist.ExternalUserId) < 1 {
		isUpdate = true
		// Use when routing is pitel_bss_conversation_
		conversationExistSecond, errTmp := repository.ConversationESRepo.GetConversationById(ctx, "", ES_INDEX_CONVERSATION, appId, conversationId)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		} else if len(conversationExistSecond.ExternalUserId) < 1 {
			log.Errorf("conversation %s not found", conversationId)
			err = errors.New("conversation " + conversationId + " not found")
			return
		}
		if err = util.ParseAnyToAny(conversationExistSecond, &conversationExist); err != nil {
			log.Error(err)
			return
		}
		conversationExist = conversationExistSecond
	}

	conversationExist.ShareInfo = &shareInfo
	conversationExist.UpdatedAt = time.Now().Format(time.RFC3339)
	conversationExist.TenantId = tenantId
	if len(connectionId) > 0 {
		*conversationExist, err = CacheConnection(ctx, connectionId, *conversationExist)
		if err != nil {
			log.Error(err)
			return
		}
	}
	if conversationExist.Labels == nil {
		conversationExist.Labels = []byte("[]")
	}
	tmpBytes, err := json.Marshal(conversationExist)
	if err != nil {
		log.Error(err)
		return
	}
	esDoc := map[string]any{}
	if err = json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return
	}

	// PROBLEM: Use for conv have not authUser and then having authUser
	// TODO: insert new conv
	if isUpdate {
		if err = repository.ESRepo.DeleteById(ctx, ES_INDEX_CONVERSATION, conversationId); err != nil {
			log.Error(err)
			return
		}
		if err = repository.ESRepo.InsertLog(ctx, tenantId, ES_INDEX_CONVERSATION, appId, conversationId, esDoc); err != nil {
			log.Error(err)
			return
		}
	} else {
		conversationQueue := model.ConversationQueue{
			DocId:        conversationExist.ConversationId,
			Conversation: *conversationExist,
		}
		if err = PublishPutConversationToChatQueue(ctx, conversationQueue); err != nil {
			log.Error(err)
			return
		}
	}

	if err = cache.RCache.Set(CONVERSATION+"_"+conversationId, conversationExist, CONVERSATION_EXPIRE); err != nil {
		log.Error(err)
		return
	}

	return
}
