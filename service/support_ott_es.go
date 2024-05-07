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

func UpSertConversation(ctx context.Context, connectionId string, data model.OttMessage) (conversation model.Conversation, isNew bool, err error) {
	newConversationId := GenerateConversationId(data.AppId, data.OaId, data.ExternalUserId)
	conversation = model.Conversation{
		TenantId:         data.TenantId,
		ConversationId:   newConversationId,
		AppId:            data.AppId,
		ConversationType: data.MessageType,
		Username:         data.Username,
		Avatar:           data.Avatar,
		OaId:             data.OaId,
		ShareInfo:        data.ShareInfo,
		ExternalUserId:   data.ExternalUserId,
		CreatedAt:        time.Now().Format(time.RFC3339),
	}
	shareInfo := data.ShareInfo

	isExisted := false
	// conversationCache := cache.RCache.Get(CONVERSATION + "_" + newConversationId)
	// if conversationCache != nil {
	// 	isExisted = true
	// 	if err := json.Unmarshal([]byte(conversationCache.(string)), &conversation); err != nil {
	// 		log.Error(err)
	// 		return conversation, isNew, err
	// 	}
	// 	if err := UpdateESAndCache(ctx, data.TenantId, data.AppId, data.ExternalUserId, connectionId, *conversation.ShareInfo); err != nil {
	// 		log.Error(err)
	// 		return conversation, isNew, err
	// 	}
	// 	return conversation, isNew, nil
	// } else {
	filter := model.ConversationFilter{
		ConversationId: []string{newConversationId},
		AppId:          []string{data.AppId},
	}
	total, conversations, err := repository.ConversationESRepo.GetConversations(ctx, "", ES_INDEX_CONVERSATION, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return conversation, isNew, err
	}
	if total > 0 {
		conversation.TenantId = (*conversations)[0].TenantId
		conversation.ConversationId = (*conversations)[0].ConversationId
		conversation.ConversationType = (*conversations)[0].ConversationType
		conversation.AppId = (*conversations)[0].AppId
		conversation.OaId = (*conversations)[0].OaId
		conversation.OaName = (*conversations)[0].OaName
		conversation.OaAvatar = (*conversations)[0].OaAvatar
		conversation.ExternalUserId = (*conversations)[0].ExternalUserId
		conversation.Username = (*conversations)[0].Username
		conversation.Username = (*conversations)[0].Username
		conversation.Avatar = (*conversations)[0].Avatar
		conversation.IsDone = (*conversations)[0].IsDone
		conversation.IsDoneBy = (*conversations)[0].IsDoneBy

		conversation.ShareInfo = shareInfo
		if len(connectionId) > 0 {
			conversation, err = CacheConnection(ctx, connectionId, conversation)
			if err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
		}

		tmpBytes, err := json.Marshal(conversation)
		if err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		esDoc := map[string]any{}
		if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		newConversationId := GenerateConversationId(conversation.AppId, conversation.OaId, conversation.ExternalUserId)
		if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, conversation.AppId, newConversationId, esDoc); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		isExisted = true
		return conversation, isNew, nil
	}
	// }

	if !isExisted {
		id, err := InsertConversation(ctx, conversation, connectionId)
		if err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		conversation.ConversationId = id
		if len(connectionId) > 0 {
			conversation, err = CacheConnection(ctx, connectionId, conversation)
			if err != nil {
				log.Error(err)
				return conversation, isNew, err
			}
		}
		if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversation, CONVERSATION_EXPIRE); err != nil {
			log.Error(err)
			return conversation, isNew, err
		}
		isNew = true
	}

	return conversation, isNew, nil
}

func InsertConversation(ctx context.Context, conversation model.Conversation, connectionId string) (id string, err error) {
	id = GenerateConversationId(conversation.AppId, conversation.OaId, conversation.ExternalUserId)
	if len(connectionId) > 0 {
		conversation, err = CacheConnection(ctx, connectionId, conversation)
		if err != nil {
			log.Error(err)
			return id, err
		}
	}
	tmpBytes, err := json.Marshal(conversation)
	if err != nil {
		log.Error(err)
		return id, err
	}

	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return id, err
	}
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX_CONVERSATION, conversation.TenantId); err != nil {
		log.Error(err)
		return id, err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, ES_INDEX_CONVERSATION, conversation.TenantId); err != nil {
			log.Error(err)
			return id, err
		}
	}
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, "", ES_INDEX_CONVERSATION, conversation.AppId, id)
	if err != nil {
		log.Error(err)
		return id, err
	} else if len(conversationExist.ExternalUserId) > 0 {
		log.Errorf("conversation %s not found", id)
		return id, errors.New("conversation " + id + " not found")
	}
	if err := repository.ESRepo.InsertLog(ctx, conversation.TenantId, ES_INDEX_CONVERSATION, conversation.AppId, id, esDoc); err != nil {
		log.Error(err)
		return id, err
	}

	return id, nil
}

/**
* Update ES and Cache
* API get conversation can get from redis, and here can caching to descrese the number of api calls to ES
 */
func UpdateESAndCache(ctx context.Context, tenantId, appId, oaId, conversationId, connectionId string, shareInfo model.ShareInfo) error {
	var isUpdate bool
	newConversationId := GenerateConversationId(appId, oaId, conversationId)
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, tenantId, ES_INDEX_CONVERSATION, appId, newConversationId)
	if err != nil {
		log.Error(err)
		return err
	} else if len(conversationExist.ExternalUserId) < 1 {
		isUpdate = true
		// Use when routing is pitel_bss_conversation_
		conversationExistSecond, err := repository.ConversationESRepo.GetConversationById(ctx, "", ES_INDEX_CONVERSATION, appId, newConversationId)
		if err != nil {
			log.Error(err)
			return err
		} else if len(conversationExistSecond.ExternalUserId) < 1 {
			log.Errorf("conversation %s not found", newConversationId)
			return errors.New("conversation " + newConversationId + " not found")
		}
		if err := util.ParseAnyToAny(conversationExistSecond, &conversationExist); err != nil {
			log.Error(err)
			return err
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
			return err
		}
	}
	tmpBytes, err := json.Marshal(conversationExist)
	if err != nil {
		log.Error(err)
		return err
	}
	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return err
	}

	// PROBLEM: Use for conv have not authUser and then having authUser
	// TODO: insert new conv
	if isUpdate {
		if err := repository.ESRepo.DeleteById(ctx, ES_INDEX_CONVERSATION, newConversationId); err != nil {
			log.Error(err)
			return err
		}
		if err := repository.ESRepo.InsertLog(ctx, tenantId, ES_INDEX_CONVERSATION, appId, newConversationId, esDoc); err != nil {
			log.Error(err)
			return err
		}
	} else {
		if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, appId, newConversationId, esDoc); err != nil {
			log.Error(err)
			return err
		}
	}

	if err := cache.RCache.Set(CONVERSATION+"_"+newConversationId, conversationExist, CONVERSATION_EXPIRE); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
