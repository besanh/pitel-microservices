package grpc

import (
	pb "github.com/tel4vn/pitel-microservices/gen/proto/chat_script"
	"github.com/tel4vn/pitel-microservices/model"
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
	}
	return
}
