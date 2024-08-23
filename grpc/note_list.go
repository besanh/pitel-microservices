package grpc

import (
	"context"
	"errors"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/note_list"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCNotesList struct{}

func NewGRPCNotesList() pb.NotesListServiceServer {
	return &GRPCNotesList{}
}

func (g *GRPCNotesList) InsertNoteInConversation(ctx context.Context, request *pb.PostNoteInConversationRequest) (result *pb.PostNoteInConversationResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ConversationNoteRequest{
		Content:        request.GetContent(),
		ConversationId: request.GetConversationId(),
		AppId:          request.GetAppId(),
		OaId:           request.GetOaId(),
	}
	if err = payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.NotesListService.InsertNoteInConversation(ctx, user, &payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PostNoteInConversationResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return
}

func (g *GRPCNotesList) UpdateNoteInConversationById(ctx context.Context, request *pb.PutNoteInConversationRequest) (result *pb.PutNoteInConversationResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ConversationNoteRequest{
		Content:        request.GetContent(),
		ConversationId: request.GetConversationId(),
		AppId:          request.GetAppId(),
		OaId:           request.GetOaId(),
	}

	if err = payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = service.NotesListService.UpdateNoteInConversationById(ctx, user, request.GetNoteId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PutNoteInConversationResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCNotesList) DeleteNoteInConversationById(ctx context.Context, request *pb.DeleteNoteInConversationRequest) (result *pb.DeleteNoteInConversationResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	if len(request.GetNoteId()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}
	payload := model.ConversationNoteRequest{
		ConversationId: request.GetConversationId(),
		AppId:          request.GetAppId(),
		OaId:           request.GetOaId(),
	}

	err = service.NotesListService.DeleteNoteInConversationById(ctx, user, request.GetNoteId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.DeleteNoteInConversationResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCNotesList) GetConversationNotesList(ctx context.Context, request *pb.GetConversationNotesListRequest) (result *pb.GetConversationNotesListResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	limit := util.ParseLimit(request.GetLimit())
	offset := util.ParseLimit(request.GetOffset())

	filter := model.ConversationNotesListFilter{
		ConversationId: request.GetConversationId(),
		AppId:          request.GetAppId(),
		OaId:           request.GetOaId(),
	}
	if len(filter.ConversationId) < 1 {
		err = errors.New("conversation id is required")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	total, data, err := service.NotesListService.GetConversationNotesList(ctx, user, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.NotesList, 0)
	for _, item := range *data {
		tmp := pb.NotesList{}
		if err = util.ParseAnyToAny(item, &tmp); err != nil {
			log.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
		tmp.CreatedAt = timestamppb.New(item.CreatedAt)
		tmp.UpdatedAt = timestamppb.New(item.UpdatedAt)

		resultData = append(resultData, &tmp)
	}
	result = &pb.GetConversationNotesListResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return
}
