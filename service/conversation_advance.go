package service

import (
	"context"
	"encoding/json"
	"errors"
	"slices"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
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
		LabelType: labelType,
		LabelName: request.LabelName,
	}
	_, chatLabelExist, err := repository.ChatLabelRepo.GetChatLabels(ctx, dbCon, filter, 1, 0)
	if err != nil {
		log.Error(err)
		return
	}
	if request.Action == "create" {
		if len(*chatLabelExist) > 0 {
			log.Error("chat label " + request.LabelName + " already exists")
			err = errors.New("chat label " + request.LabelName + " already exists")
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
		LabelColor:  "",
		LabelStatus: true,
		CreatedBy:   authUser.UserId,
	}

	var externalLabelId string

	// TODO: create label zalo
	if labelType == "zalo" {
		if err = s.handleLabelZalo(ctx, labelType, request); err != nil {
			return
		}
		if request.Action == "create" {
			labelId = chatLabel.GetId()
		} else if request.Action == "update" || request.Action == "delete" {
			labelId = (*chatLabelExist)[0].GetId()
		}
	} else if labelType == "facebook" {
		externalLabelId, err = s.handleLabelFacebook(ctx, dbCon, labelType, chatLabel, request)
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
		labelId = chatLabel.GetId()
	} else if request.Action == "update" {
		// TODO: we don't need to do anything
	} else if request.Action == "delete" {
		// TODO: we don't need to do anything
	}

	// TODO: update label for conversation => use queue
	if err = s.putConversation(ctx, authUser, labelId, labelType, request); err != nil {
		if request.Action == "create" {
			if err = repository.ChatLabelRepo.Delete(ctx, dbCon, chatLabel.GetId()); err != nil {
				log.Error(err)
				return
			}
		}
		return
	}

	return
}

func (s *Conversation) handleLabelZalo(ctx context.Context, labelType string, request model.ConversationLabelRequest) (err error) {
	zaloRequest := model.ChatExternalLabelRequest{
		AppId:          request.AppId,
		OaId:           request.OaId,
		ExternalUserId: request.ExternalUserId,
		TagName:        request.LabelName,
	}
	var externalUrl string
	if request.Action == "create" {
		externalUrl = "create-label-customer"
	} else if request.Action == "update" {
		// TODO: we don't need to do anything
	} else if request.Action == "delete" {
		externalUrl = "remove-label-customer"
	}

	// TODO: because zalo not return id so we don't need to use it for updating label
	if slices.Contains([]string{"create", "delete"}, request.Action) {
		_, errTmp := RequestOttLabel(ctx, labelType, externalUrl, zaloRequest)
		if errTmp != nil {
			log.Error(errTmp)
			err = errTmp
			return
		}
	}

	return
}

func (s *Conversation) handleLabelFacebook(ctx context.Context, dbCon sqlclient.ISqlClientConn, labelType string, chatLabel model.ChatLabel, request model.ConversationLabelRequest) (externalLabelId string, err error) {
	// TODO: if label does not exist, create new label, then associate to conversation)
	facebookRequest := model.ChatExternalLabelRequest{
		AppId:          request.AppId,
		OaId:           request.OaId,
		ExternalUserId: request.ExternalUserId,
		LabelId:        chatLabel.GetId(),
		TagName:        request.LabelName,
	}
	externalUrl := "create-label"
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

	return
}

func (s *Conversation) putConversation(ctx context.Context, authUser *model.AuthUser, labelId, labelType string, request model.ConversationLabelRequest) (err error) {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, request.AppId, request.ConversationId)
	if err != nil {
		log.Error(err)
		return
	} else if len(conversationExist.ConversationId) < 1 {
		log.Errorf("conversation %s not found", request.ConversationId)
		return
	}

	objmap := []any{}
	labalesExist := []any{}
	if err = json.Unmarshal([]byte(conversationExist.Label), &labalesExist); err != nil {
		log.Error(err)
		return
	}

	// TODO: because zalo only assign one label for one conversation
	if labelType == "facebook" {
		for _, item := range labalesExist {
			tmp := map[string]string{}
			if err = util.ParseAnyToAny(item, &tmp); err != nil {
				log.Error(err)
				continue
			}
			if request.Action == "delete" && tmp["label_id"] == labelId {
				continue
			}
			objmap = append(objmap, map[string]any{
				"label_id": tmp["label_id"],
			})
		}
	}

	if request.Action == "create" || request.Action == "update" {
		objmap = append(objmap, map[string]any{
			"label_id": labelId,
		})
	}

	result, err := json.Marshal(objmap)
	if err != nil {
		log.Error(err)
	}
	conversationExist.Label = result

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
	if err = repository.ESRepo.UpdateDocById(ctx, ES_INDEX_CONVERSATION, request.AppId, request.ConversationId, esDoc); err != nil {
		log.Error(err)
		return
	}
	return
}
