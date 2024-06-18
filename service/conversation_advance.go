package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

/**
* Create in internal system, then create in external
 */
func (s *Conversation) PutLabelToConversation(ctx context.Context, authUser *model.AuthUser, labelType string, request model.ConversationLabelRequest) (labelId string, err error) {
	dbCon, err := HandleGetDBConSource(authUser)
	if err != nil {
		log.Error(err)
		return
	}
	filter := model.ChatLabelFilter{
		TenantId:  authUser.TenantId,
		AppId:     request.AppId,
		OaId:      request.OaId,
		LabelName: request.LabelName,
	}
	_, labelExist, err := repository.ChatLabelRepo.GetChatLabels(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	} else if len(*labelExist) > 0 {
		log.Error("chat label " + request.LabelName + " already exists")
		err = errors.New("chat label " + request.LabelName + " already exists")
		return
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
		LabelColor:  "",
		LabelStatus: true,
		CreatedBy:   authUser.UserId,
	}
	if err = repository.ChatLabelRepo.Insert(ctx, dbCon, chatLabel); err != nil {
		log.Error(err)
		return
	}

	var externalUrl string
	var externalLabelId string

	// TODO: create label zalo
	if labelType == "zalo" {
		zaloRequest := model.ChatExternalLabelRequest{
			AppId:          request.AppId,
			OaId:           request.OaId,
			ExternalUserId: request.ExternalUserId,
			TagName:        request.LabelName,
		}
		externalUrl = "create-label-customer"

		// TODO: because zalo not return id so we don't need to use it for updating label
		_, errTmp := RequestOttLabel(ctx, labelType, externalUrl, zaloRequest)
		if errTmp != nil {
			log.Error(errTmp)
			if err = repository.ChatLabelRepo.Delete(ctx, dbCon, chatLabel.GetId()); err != nil {
				log.Error(err)
				return
			}
			err = errTmp
			return
		}
	} else if labelType == "facebook" {
		// TODO: if label does not exist, create new label, then associate to conversation)
		facebookRequest := model.ChatExternalLabelRequest{
			AppId:          request.AppId,
			OaId:           request.OaId,
			ExternalUserId: request.ExternalUserId,
			LabelId:        chatLabel.GetId(),
			TagName:        request.LabelName,
		}
		externalUrl = "create-label"
		externalCreateLabelResponse, errTmp := RequestOttLabel(ctx, labelType, externalUrl, facebookRequest)
		if errTmp != nil {
			log.Error(errTmp)
			if err = repository.ChatLabelRepo.Delete(ctx, dbCon, chatLabel.GetId()); err != nil {
				log.Error(err)
				return
			}
			return
		}
		externalLabelId = externalCreateLabelResponse.Id

		// TODO: associating a label to oa
		facebookAssociateRequest := model.ChatExternalLabelRequest{
			AppId:          request.AppId,
			OaId:           request.OaId,
			ExternalUserId: externalLabelId,
			LabelId:        chatLabel.GetId(),
			TagName:        request.LabelName,
		}
		externalUrl = "associate-label"
		externalAssociateLabelResponse, errTmp := RequestOttLabel(ctx, labelType, externalUrl, facebookAssociateRequest)
		if errTmp != nil {
			log.Error(errTmp)
			if err = repository.ChatLabelRepo.Delete(ctx, dbCon, chatLabel.GetId()); err != nil {
				log.Error(err)
				return
			}
			err = errTmp
			return
		}
		if len(externalAssociateLabelResponse.Id) > 0 {
			externalLabelId = externalAssociateLabelResponse.Id
		} else {
			err = errors.New("external label id not found")
			log.Error("external label id not found")
			return
		}

		// TODO: update label
		chatLabel.ExternalLabelId = externalLabelId
	}

	// TODO: update label
	chatLabel.UpdatedAt = time.Now()
	if err = repository.ChatLabelRepo.Update(ctx, dbCon, chatLabel); err != nil {
		log.Error(err)
		return
	}

	// TODO: update label for conversation
	conversationId := GenerateConversationId(request.AppId, request.OaId, request.ConversationId)
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, request.AppId, conversationId)
	if err != nil {
		log.Error(err)
		return
	} else if len(conversationExist.ConversationId) < 1 {
		log.Errorf("conversation %s not found", conversationId)
		return
	}
	conversationExist.Label = append(conversationExist.Label, chatLabel.GetId())
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
	if err = repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, request.AppId, conversationId, esDoc); err != nil {
		log.Error(err)
		return
	}

	return
}
