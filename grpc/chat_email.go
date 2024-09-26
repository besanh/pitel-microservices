package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/tel4vn/pitel-microservices/common/log"
	"github.com/tel4vn/pitel-microservices/common/response"
	"github.com/tel4vn/pitel-microservices/common/util"
	pb "github.com/tel4vn/pitel-microservices/gen/proto/chat_email"
	"github.com/tel4vn/pitel-microservices/middleware/auth"
	"github.com/tel4vn/pitel-microservices/model"
	"github.com/tel4vn/pitel-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCChatEmail struct{}

func NewGRPCChatEmail() pb.ChatEmailServiceServer {
	return &GRPCChatEmail{}
}

func (g *GRPCChatEmail) InsertChatEmail(ctx context.Context, request *pb.PostChatEmailRequest) (*pb.PostChatEmailResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatEmailRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	id, err := service.ChatEmailService.InsertChatEmail(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PostChatEmailResponse{
		Code:    "OK",
		Message: "ok",
		Id:      id,
	}
	return result, nil
}

func (g *GRPCChatEmail) GetChatEmails(ctx context.Context, request *pb.GetChatEmailsRequest) (*pb.GetChatEmailsResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatEmailFilter{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	limit, offset := request.GetLimit(), request.GetOffset()
	statusTmp := request.GetStatus()
	var scriptStatus sql.NullBool
	if len(statusTmp) > 0 {
		tmp, _ := strconv.ParseBool(statusTmp)
		scriptStatus.Valid = true
		scriptStatus.Bool = tmp
	}
	payload.Status = scriptStatus

	total, data, err := service.ChatEmailService.GetChatEmails(ctx, user, payload, int(limit), int(offset))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ChatEmailCustomData, 0)
	if len(*data) > 0 {
		for _, item := range *data {
			var tmp pb.ChatEmailCustomData
			tmp.CreatedAt = &timestamppb.Timestamp{
				Seconds: item.CreatedAt.Unix(),
			}
			tmp.UpdatedAt = &timestamppb.Timestamp{
				Seconds: item.UpdatedAt.Unix(),
			}
			if err = util.ParseAnyToAny(item, &tmp); err != nil {
				log.Error(err)
				result := &pb.GetChatEmailsResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return result, nil
			}
			resultData = append(resultData, &tmp)
		}
	}

	result := &pb.GetChatEmailsResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   limit,
		Offset:  offset,
	}
	return result, nil
}

func (g *GRPCChatEmail) GetChatEmailById(ctx context.Context, request *pb.GetEmailByIdRequest) (*pb.GetEmailByIdResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	data, err := service.ChatEmailService.GetChatEmailById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.ChatEmailData{}
	tmp.CreatedAt = &timestamppb.Timestamp{
		Seconds: data.CreatedAt.Unix(),
	}
	tmp.UpdatedAt = &timestamppb.Timestamp{
		Seconds: data.UpdatedAt.Unix(),
	}
	if err = util.ParseAnyToAny(data, tmp); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.GetEmailByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
	}
	return result, nil
}

func (g *GRPCChatEmail) UpdateChatEmailById(ctx context.Context, request *pb.PutChatEmailRequest) (*pb.PutChatEmailResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	payload := model.ChatEmailRequest{}
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ChatEmailService.UpdateChatEmailById(ctx, user, request.GetId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutChatEmailResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCChatEmail) DeleteChatEmailById(ctx context.Context, request *pb.DeleteChatEmailRequest) (*pb.DeleteChatEmailResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	err := service.ChatEmailService.DeleteChatEmailById(ctx, user, request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.DeleteChatEmailResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}
