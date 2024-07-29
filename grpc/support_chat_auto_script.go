package grpc

import (
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_auto_script"
	"github.com/tel4vn/fins-microservices/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertChatAutoScriptToPbChatAutoScript(data model.ChatAutoScriptView) (result *pb.ChatAutoScriptData, err error) {
	result = &pb.ChatAutoScriptData{
		Id:                 data.Id,
		CreatedAt:          timestamppb.New(data.CreatedAt),
		UpdatedAt:          timestamppb.New(data.UpdatedAt),
		TenantId:           data.TenantId,
		ScriptName:         data.ScriptName,
		Status:             data.Status,
		ConnectionId:       data.ConnectionId,
		ConnectionApp:      nil,
		Channel:            data.Channel,
		CreatedBy:          data.CreatedBy,
		UpdatedBy:          data.UpdatedBy,
		TriggerEvent:       data.TriggerEvent,
		TriggerKeywords:    nil,
		ChatScriptLink:     nil,
		SendMessageActions: nil,
		ChatLabelLink:      make([]*pb.ChatLabelLinkDataType, 0),
		ActionScript:       nil,
	}
	if data.ConnectionApp != nil {
		result.ConnectionApp = &pb.ChatConnectionAppData{
			Id:                data.ConnectionApp.Id,
			TenantId:          data.ConnectionApp.TenantId,
			CreatedAt:         timestamppb.New(data.ConnectionApp.CreatedAt),
			UpdatedAt:         timestamppb.New(data.ConnectionApp.UpdatedAt),
			ConnectionName:    data.ConnectionApp.ConnectionName,
			ConnectionType:    data.ConnectionApp.ConnectionType,
			ChatAppId:         data.ConnectionApp.ChatAppId,
			Status:            data.ConnectionApp.Status,
			ConnectionQueueId: data.ConnectionApp.ConnectionQueueId,
		}
		if err = util.ParseAnyToAny(data.ConnectionApp.OaInfo, &result.ConnectionApp.OaInfo); err != nil {
			return
		}
	}
	if err = util.ParseAnyToAny(data.TriggerKeywords, &result.TriggerKeywords); err != nil {
		return
	}
	if err = util.ParseAnyToAny(data.ChatScriptLink, &result.ChatScriptLink); err != nil {
		return
	}
	if err = util.ParseAnyToAny(data.SendMessageActions, &result.SendMessageActions); err != nil {
		return
	}
	for _, label := range data.ChatLabelLink {
		if label != nil {
			tmp := &pb.ChatLabelLinkDataType{
				ChatAutoScriptId: label.ChatAutoScriptId,
				ChatLabelId:      label.ChatLabelId,
				ActionType:       label.ActionType,
				Order:            int32(label.Order),
				ChatAutoScript:   nil,
				ChatLabel:        nil,
				CreatedAt:        timestamppb.New(label.CreatedAt),
				UpdatedAt:        timestamppb.New(label.UpdatedAt),
			}
			if label.ChatLabel != nil {
				if err = util.ParseAnyToAny(label.ChatLabel, &tmp.ChatLabel); err != nil {
					return
				}
			}
			result.ChatLabelLink = append(result.ChatLabelLink, tmp)
		}
	}
	if data.ActionScript != nil {
		if err = util.ParseAnyToAny(data.ActionScript, &result.ActionScript); err != nil {
			return
		}
	}
	return
}
