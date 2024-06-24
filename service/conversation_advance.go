package service

import (
	"context"
	"database/sql"
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
			externalLabelId = chatLabel.GetId()
		}
	} else if request.Action == "delete" {
		if labelType == "facebook" {
			externalLabelId = request.ExternalLabelId
		}
	}

	// TODO: update label for conversation => use queue
	if err = putConversation(ctx, authUser, externalLabelId, labelType, request); err != nil {
		if request.Action == "create" {
			if err = repository.ChatLabelRepo.Delete(ctx, dbCon, chatLabel.GetId()); err != nil {
				log.Error(err)
				return
			}
		}
		return
	}

	labelId = chatLabel.GetId()

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

func putConversation(ctx context.Context, authUser *model.AuthUser, labelId, labelType string, request model.ConversationLabelRequest) (err error) {
	conversationExist, err := repository.ConversationESRepo.GetConversationById(ctx, authUser.TenantId, ES_INDEX_CONVERSATION, request.AppId, request.ConversationId)
	if err != nil {
		log.Error(err)
		return
	} else if len(conversationExist.ConversationId) < 1 {
		log.Errorf("conversation %s not found", request.ConversationId)
		return
	}

	objmap := []any{}
	labelsExist := []any{}
	if err = json.Unmarshal([]byte(conversationExist.Label), &labelsExist); err != nil {
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
			if request.Action == "delete" && tmp["label_id"] == labelId {
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

	if request.Action == "create" || request.Action == "update" {
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
	}
	log.Info("label_id: ", labelId)

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

func checkItemExist(objmap []any, tmp map[string]string) (isExist bool) {
	for _, item := range objmap {
		if item.(map[string]any)["label_id"] == tmp["label_id"] {
			return true
		}
	}
	return
}
