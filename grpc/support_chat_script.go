package grpc

import (
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_script"
	"github.com/tel4vn/fins-microservices/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertChatScriptToPbChatScript(data model.ChatScriptView) (result *pb.ChatScriptData, err error) {
	result = &pb.ChatScriptData{
		Id:            data.Id,
		CreatedAt:     timestamppb.New(data.CreatedAt),
		UpdatedAt:     timestamppb.New(data.UpdatedAt),
		TenantId:      data.TenantId,
		ScriptName:    data.ScriptName,
		Status:        data.Status,
		ScriptType:    data.ScriptType,
		FileUrl:       data.FileUrl,
		Channel:       data.Channel,
		Content:       data.Content,
		CreatedBy:     data.CreatedBy,
		UpdatedBy:     data.UpdatedBy,
		OtherScriptId: data.OtherScriptId,
		ConnectionId:  data.ConnectionId,
		ConnectionApp: nil,
	}
	if data.ConnectionApp != nil {
		result.ConnectionApp = &pb.ChatConnectionAppData{
			Id:                data.ConnectionApp.Id,
			TenantId:          data.ConnectionApp.Id,
			CreatedAt:         timestamppb.New(data.ConnectionApp.CreatedAt),
			UpdatedAt:         timestamppb.New(data.ConnectionApp.UpdatedAt),
			ConnectionName:    data.ConnectionApp.Id,
			ConnectionType:    data.ConnectionApp.Id,
			ChatAppId:         data.ConnectionApp.Id,
			Status:            data.ConnectionApp.Id,
			ConnectionQueueId: data.ConnectionApp.Id,
		}
		if err = util.ParseAnyToAny(data.ConnectionApp.OaInfo, &result.ConnectionApp.OaInfo); err != nil {
			return
		}
	}
	return
}
