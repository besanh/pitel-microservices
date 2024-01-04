package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
	"github.com/tel4vn/fins-microservices/common/variables"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"golang.org/x/exp/slices"
)

type (
	IOttMessage interface {
		GetOttMessage(ctx context.Context, data model.GetOttMessage) error
	}
	OttMessage struct{}
)

func NewOttMessage() IOttMessage {
	return &OttMessage{}
}

func (s *OttMessage) GetOttMessage(ctx context.Context, data model.GetOttMessage) error {
	docId := uuid.NewString()
	timestamp := time.Unix(0, data.Timestamp*int64(time.Millisecond))
	message := model.Message{
		Id:            docId,
		ParentMsgId:   "",
		MsgId:         data.MsgId,
		MessageType:   data.Type,
		EventName:     data.EventName,
		Direction:     variables.DIRECTION["receive"],
		AppId:         data.AppId,
		OaId:          data.OaId,
		UserIdByApp:   data.UserIdByApp,
		UserId:        data.UserId,
		Username:      data.Username,
		Avatar:        data.Avatar,
		SendTime:      timestamp,
		SendTimestamp: data.Timestamp,
		Content:       data.Text,
	}
	if slices.Contains[[]string](variables.EVENT_READ_MESSAGE, data.EventName) {
		timestamp := time.Unix(0, data.Timestamp*int64(time.Millisecond))
		message.ReadTime = timestamp
		message.ReadTimestamp = data.Timestamp
	}
	if data.Attachments != nil {
		for _, val := range data.Attachments {
			var attachmentFile model.OttPayloadFile
			var attachmentMedia model.OttPayloadMedia
			if val.AttType == variables.ATTACHMENT_TYPE["file"] {
				if err := util.ParseAnyToAny(val.Payload, attachmentFile); err != nil {
					log.Error(err)
					return err
				}
			} else {
				if err := util.ParseAnyToAny(val.Payload, attachmentMedia); err != nil {
					log.Error(err)
					return err
				}
			}
			message.Attachments = append(message.Attachments, &model.Attachments{
				Id:             uuid.NewString(),
				MsgId:          docId,
				AttachmentType: val.AttType,
				AttachmentsDetail: &model.AttachmentsDetail{
					AttachmentFile:  &attachmentFile,
					AttachmentMedia: &attachmentMedia,
				},
				SendTime:      timestamp,
				SendTimestamp: data.Timestamp,
			})
		}
	}

	//  TODO: add rabbitmq message
	tmpBytes, err := json.Marshal(message)
	if err != nil {
		log.Error(err)
		return err
	}

	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return err
	}
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX, data.AppId); err != nil {
		log.Error(err)
		return err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, ES_INDEX, data.AppId); err != nil {
			log.Error(err)
			return err
		}
	}

	if err := HandlePushRMQ(ctx, ES_INDEX, docId, message, tmpBytes); err != nil {
		log.Error(err)
		return err
	}
	// TODO: delivery queue
	// if err

	return nil
}
