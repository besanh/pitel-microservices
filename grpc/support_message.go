package grpc

import (
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/message"
	"github.com/tel4vn/fins-microservices/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertMessageToPbMessage(data model.Message) (result *pb.Message, err error) {
	result = &pb.Message{
		TenantId:            data.TenantId,
		ParentMessageId:     data.ParentMessageId,
		Id:                  data.Id,
		ConversationId:      data.ConversationId,
		ParentExternalMsgId: data.ParentExternalMsgId,
		ExternalMsgId:       data.ExternalMsgId,
		MessageType:         data.MessageType,
		EventName:           data.EventName,
		Direction:           data.Direction,
		AppId:               data.AppId,
		OaId:                data.OaId,
		UserIdByApp:         data.UserIdByApp,
		ExternalUserId:      data.ExternalUserId,
		UserAppName:         data.UserAppname,
		Avatar:              data.Avatar,
		SupporterId:         data.SupporterId,
		SupporterName:       data.SupporterName,
		SendTime:            timestamppb.New(data.SendTime),
		SendTimestamp:       data.SendTimestamp,
		Content:             data.Content,
		IsRead:              data.IsRead,
		ReadTime:            timestamppb.New(data.ReadTime),
		ReadTimestamp:       data.ReadTimestamp,
		ReadBy:              data.ReadBy,
		Attachments:         nil,
		CreatedAt:           timestamppb.New(data.CreatedAt),
		UpdatedAt:           timestamppb.New(data.UpdatedAt),
		ShareInfo:           nil,
		IsEcho:              data.IsEcho,
	}
	if len(data.Attachments) > 0 {
		if err = util.ParseAnyToAny(data.Attachments, &result.Attachments); err != nil {
			return
		}
	}
	if data.ShareInfo != nil {
		if err = util.ParseAnyToAny(data.ShareInfo, &result.ShareInfo); err != nil {
			return
		}
	}

	return
}
