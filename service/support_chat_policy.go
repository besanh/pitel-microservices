package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func GeneratePolicySettingKeyId(tenantId, connectionType string) string {
	return CHAT_POLICY_SETTING + "_" + tenantId + "_" + connectionType
}

/*
 * Check if it's still inside chat-able window. If we send message, 3rd parties are still going to reject them
 * This func will throw an error when that happens
 */
func CheckOutOfChatWindowTime(ctx context.Context, tenantId, connectionType, lastMessageTimestamp string) (err error) {
	chatWindowDuration := 0
	// use from env configs
	switch connectionType {
	case "zalo":
		chatWindowDuration = ZALO_POLICY_CHAT_WINDOW
	case "facebook":
		chatWindowDuration = FACEBOOK_POLICY_CHAT_WINDOW
	default:
		return errors.New("not supported connection " + connectionType)
	}

	// check if conversation's still inside chat-able window
	policySetting := model.ChatPolicySetting{}
	key := GeneratePolicySettingKeyId(tenantId, connectionType)
	policySettingCache := cache.RCache.Get(key)
	if policySettingCache != nil {
		if err = json.Unmarshal([]byte(policySettingCache.(string)), &policySetting); err != nil {
			log.Error(err)
			return
		}

		chatWindowDuration = policySetting.ChatWindowTime
	} else {
		filter := model.ChatPolicyFilter{
			TenantId:       tenantId,
			ConnectionType: connectionType,
		}
		total, policySettings, errTmp := repository.ChatPolicySettingRepo.GetChatPolicySettings(ctx, repository.DBConn, filter, 0, 0)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		}
		if total > 0 {
			// use setting from admin configs
			chatWindowDuration = (*policySettings)[0].ChatWindowTime

			// set cache
			if err = cache.RCache.Set(key, (*policySettings)[0], CHAT_POLICY_SETTING_EXPIRE); err != nil {
				log.Error(err)
				return
			}
		}
	}

	if len(lastMessageTimestamp) > 0 {
		currentTime := time.Now()
		lastMessageAt, errTmp := time.Parse(time.RFC3339, lastMessageTimestamp)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			return
		}
		if chatWindowDuration < 0 {
			chatWindowDuration = 0
		}
		windowEndTime := lastMessageAt.Add(time.Duration(chatWindowDuration) * time.Second)

		if currentTime.After(windowEndTime) {
			err = errors.New("out of chat window time")
			return
		}
	}

	return
}
