package service

import (
	"context"
	"encoding/json"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IProfile interface {
		GetUpdateProfileByUserId(ctx context.Context, authUser *model.AuthUser, request model.ProfileRequest) (int, any)
	}
	Profile struct{}
)

func NewProfile() IProfile {
	return &Profile{}
}

func (s *Profile) GetUpdateProfileByUserId(ctx context.Context, authUser *model.AuthUser, request model.ProfileRequest) (int, any) {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, request.AppId, request.ConversationId)
	if err != nil {
		log.Error(err)
		return response.ServiceUnavailableMsg(err.Error())
	} else if len(conversationExist.ConversationId) < 1 {
		log.Errorf("conversation %s not found", request.ConversationId)
		return response.ServiceUnavailableMsg("conversation " + request.ConversationId + " not found")
	}

	if len(conversationExist.ShareInfo.Fullname) < 1 || len(conversationExist.ShareInfo.PhoneNumber) < 1 {
		if request.ProfileType == "zalo" {
			res, err := GetProfile(ctx, request.AppId, request.OaId, request.UserId)
			if err != nil {
				log.Error(err)
				return response.ServiceUnavailableMsg(err.Error())
			}
			if len(res.Data.ShareInfo.Name) > 0 || len(res.Data.ShareInfo.Phone) > 0 {
				conversationExist.ShareInfo.Address = res.Data.ShareInfo.Address
				conversationExist.ShareInfo.Fullname = res.Data.ShareInfo.Name
				conversationExist.ShareInfo.PhoneNumber = res.Data.ShareInfo.Phone
				conversationExist.ShareInfo.City = res.Data.ShareInfo.City
				conversationExist.ShareInfo.District = res.Data.ShareInfo.District

				conversationExist.Username = res.Data.DisplayName
				conversationExist.Avatar = res.Data.Avatar

				tmpBytes, err := json.Marshal(conversationExist)
				if err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
				esDoc := map[string]any{}
				if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}

				if err := repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, conversationExist.AppId, conversationExist.ConversationId, esDoc); err != nil {
					log.Error(err)
					return response.ServiceUnavailableMsg(err.Error())
				}
			}
		}
	}
	return response.OK(conversationExist)
}
