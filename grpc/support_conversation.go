package grpc

import (
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/conversation"
	"github.com/tel4vn/fins-microservices/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertConversationViewToPbConversationView(data *model.ConversationView) (result *pb.ConversationView, err error) {
	result = &pb.ConversationView{
		TenantId:               data.TenantId,
		ConversationId:         data.ConversationId,
		ConversationType:       data.ConversationType,
		AppId:                  data.AppId,
		OaId:                   data.OaId,
		OaName:                 data.OaName,
		OaAvatar:               data.OaAvatar,
		ShareInfo:              nil,
		ExternalUserId:         data.ExternalUserId,
		Username:               data.Username,
		Avatar:                 data.Avatar,
		Major:                  data.Major,
		Following:              data.Following,
		Labels:                 make([]*pb.ChatLabel, 0),
		IsDone:                 data.IsDone,
		IsDoneAt:               data.IsDoneAt,
		IsDoneBy:               data.IsDoneBy,
		CreatedAt:              data.CreatedAt,
		UpdatedAt:              data.UpdatedAt,
		TotalUnread:            int32(data.TotalUnRead),
		LatestMessageContent:   data.LatestMessageContent,
		LatestMessageDirection: data.LatestMessageDirection,
		ExternalConversationId: data.ExternalConversationId,
		NotesList:              nil,
	}
	if data.ShareInfo != nil {
		if err = util.ParseAnyToAny(data.ShareInfo, &result.ShareInfo); err != nil {
			return
		}
	}
	if err = util.ParseStringToAny(string(data.Labels), &result.Labels); err != nil {
		return
	}
	if result.NotesList, err = convertNotesListToPbNotesList(data.NotesList); err != nil {
		return
	}
	return
}

func convertConversationToPbConversation(data *model.Conversation) (result *pb.Conversation, err error) {
	result = &pb.Conversation{
		TenantId:               data.TenantId,
		ConversationId:         data.ConversationId,
		ConversationType:       data.ConversationType,
		AppId:                  data.AppId,
		OaId:                   data.OaId,
		OaName:                 data.OaName,
		OaAvatar:               data.OaAvatar,
		ShareInfo:              nil,
		ExternalUserId:         data.ExternalUserId,
		Username:               data.Username,
		Avatar:                 data.Avatar,
		Major:                  data.Major,
		Following:              data.Following,
		Labels:                 make([]*pb.ChatLabel, 0),
		IsDone:                 data.IsDone,
		IsDoneAt:               timestamppb.New(data.IsDoneAt),
		IsDoneBy:               data.IsDoneBy,
		CreatedAt:              data.CreatedAt,
		UpdatedAt:              data.UpdatedAt,
		ExternalConversationId: data.ExternalConversationId,
		NotesList:              nil,
	}
	if data.ShareInfo != nil {
		if err = util.ParseAnyToAny(data.ShareInfo, &result.ShareInfo); err != nil {
			return
		}
	}
	if err = util.ParseStringToAny(string(data.Labels), &result.Labels); err != nil {
		return
	}
	if result.NotesList, err = convertNotesListToPbNotesList(data.NotesList); err != nil {
		return
	}
	return
}

func convertNotesListToPbNotesList(data *[]model.NotesList) (result []*pb.NotesList, err error) {
	if data == nil {
		return
	}
	for _, note := range *data {
		tmp := &pb.NotesList{
			Id:        note.Id,
			Content:   note.Content,
			CreatedAt: timestamppb.New(note.CreatedAt),
			UpdatedAt: timestamppb.New(note.UpdatedAt),
		}
		result = append(result, tmp)
	}
	return
}
