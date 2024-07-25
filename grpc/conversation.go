package grpc

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/conversation"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCConversation struct{}

func NewGRPCConversation() pb.ConversationServiceServer {
	return &GRPCConversation{}
}

func (g *GRPCConversation) GetConversations(ctx context.Context, request *pb.GetConversationsRequest) (result *pb.GetConversationsResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	isDone := sql.NullBool{}
	if len(request.GetIsDone()) > 0 {
		isDone.Valid = true
		isDone.Bool, _ = strconv.ParseBool(request.GetIsDone())
	}
	major := sql.NullBool{}
	if len(request.GetMajor()) > 0 {
		major.Valid = true
		major.Bool, _ = strconv.ParseBool(request.GetMajor())
	}
	following := sql.NullBool{}
	if len(request.GetFollowing()) > 0 {
		following.Valid = true
		following.Bool, _ = strconv.ParseBool(request.GetFollowing())
	}

	filter := model.ConversationFilter{
		AppId:          util.ParseQueryArray(request.GetAppId()),
		ConversationId: util.ParseQueryArray(request.GetConversationId()),
		Username:       request.GetUsername(),
		PhoneNumber:    request.GetPhoneNumber(),
		Email:          request.GetEmail(),
		IsDone:         isDone,
		Major:          major,
		Following:      following,
	}

	total, responseData, err := service.ConversationService.GetConversations(ctx, user, filter, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ConversationCustomView, 0)
	for _, item := range responseData {
		var tmp pb.ConversationCustomView
		if err = util.ParseAnyToAny(item, &tmp); err != nil {
			log.Error(err)
			result = &pb.GetConversationsResponse{
				Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
				Message: err.Error(),
			}
			return result, nil
		}
		resultData = append(resultData, &tmp)
	}

	result = &pb.GetConversationsResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return
}

func (g *GRPCConversation) GetConversationsWithScrollAPI(ctx context.Context, request *pb.GetConversationsWithScrollAPIRequest) (result *pb.GetConversationsWithScrollAPIResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	limit := util.ParseLimit(request.GetLimit())

	isDone := sql.NullBool{}
	if len(request.GetIsDone()) > 0 {
		isDone.Valid = true
		isDone.Bool, _ = strconv.ParseBool(request.GetIsDone())
	}
	major := sql.NullBool{}
	if len(request.GetMajor()) > 0 {
		major.Valid = true
		major.Bool, _ = strconv.ParseBool(request.GetMajor())
	}
	following := sql.NullBool{}
	if len(request.GetFollowing()) > 0 {
		following.Valid = true
		following.Bool, _ = strconv.ParseBool(request.GetFollowing())
	}

	filter := model.ConversationFilter{
		AppId:          util.ParseQueryArray(request.GetAppId()),
		ConversationId: util.ParseQueryArray(request.GetConversationId()),
		Username:       request.GetUsername(),
		PhoneNumber:    request.GetPhoneNumber(),
		Email:          request.GetEmail(),
		IsDone:         isDone,
		Major:          major,
		Following:      following,
	}

	code, responseData := service.ConversationService.GetConversationsWithScrollAPI(ctx, user, filter, limit, request.GetScrollId())
	result, err = parseResponseDataOfGetConversationsWithScrollAPI(code, responseData, limit)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return
}

func (g *GRPCConversation) GetConversationsByManager(ctx context.Context, request *pb.GetConversationsByManagerRequest) (result *pb.GetConversationsByManagerResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	isDone := sql.NullBool{}
	if len(request.GetIsDone()) > 0 {
		isDone.Valid = true
		isDone.Bool, _ = strconv.ParseBool(request.GetIsDone())
	}
	major := sql.NullBool{}
	if len(request.GetMajor()) > 0 {
		major.Valid = true
		major.Bool, _ = strconv.ParseBool(request.GetMajor())
	}
	following := sql.NullBool{}
	if len(request.GetFollowing()) > 0 {
		following.Valid = true
		following.Bool, _ = strconv.ParseBool(request.GetFollowing())
	}

	filter := model.ConversationFilter{
		AppId:          util.ParseQueryArray(request.GetAppId()),
		ConversationId: util.ParseQueryArray(request.GetConversationId()),
		Username:       request.GetUsername(),
		PhoneNumber:    request.GetPhoneNumber(),
		Email:          request.GetEmail(),
		IsDone:         isDone,
		Major:          major,
		Following:      following,
	}

	code, responseData := service.ConversationService.GetConversationsByHighLevel(ctx, user, filter, limit, offset)
	result, err = parseResponseDataOfGetConversationsByManager(code, responseData, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return
}

func (g *GRPCConversation) GetConversationsByManagerWithScrollAPI(ctx context.Context, request *pb.GetConversationsByManagerWithScrollAPIRequest) (result *pb.GetConversationsByManagerWithScrollAPIResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	limit := util.ParseLimit(request.GetLimit())

	isDone := sql.NullBool{}
	if len(request.GetIsDone()) > 0 {
		isDone.Valid = true
		isDone.Bool, _ = strconv.ParseBool(request.GetIsDone())
	}
	major := sql.NullBool{}
	if len(request.GetMajor()) > 0 {
		major.Valid = true
		major.Bool, _ = strconv.ParseBool(request.GetMajor())
	}
	following := sql.NullBool{}
	if len(request.GetFollowing()) > 0 {
		following.Valid = true
		following.Bool, _ = strconv.ParseBool(request.GetFollowing())
	}

	filter := model.ConversationFilter{
		AppId:          util.ParseQueryArray(request.GetAppId()),
		ConversationId: util.ParseQueryArray(request.GetConversationId()),
		Username:       request.GetUsername(),
		PhoneNumber:    request.GetPhoneNumber(),
		Email:          request.GetEmail(),
		IsDone:         isDone,
		Major:          major,
		Following:      following,
	}

	code, responseData := service.ConversationService.GetConversationsByHighLevelWithScrollAPI(ctx, user, filter, limit, request.GetScrollId())
	result, err = parseResponseDataOfGetConversationsByManagerWithScrollAPI(code, responseData, limit)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return
}

func (g *GRPCConversation) GetConversationById(ctx context.Context, request *pb.GetConversationByIdRequest) (result *pb.GetConversationByIdResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	conversationId := request.GetId()
	if len(conversationId) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, response.ERR_GET_FAILED)
	}
	code, responseData := service.ConversationService.GetConversationById(ctx, user, request.GetAppId(), conversationId)
	result, err = parseResponseDataOfGetConversationById(code, responseData)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return
}

func (g *GRPCConversation) UpdateConversation(ctx context.Context, request *pb.PutConversationRequest) (result *pb.PutConversationResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	if len(request.GetAppId()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "app_id is required")
	}
	if len(request.GetOaId()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "oa_id is required")
	}
	if len(request.GetId()) < 1 {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	var payload model.ShareInfo
	if err = util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	code, _ := service.ConversationService.UpdateConversationById(ctx, user, request.GetAppId(), request.GetOaId(), request.GetId(), payload)
	if code != http.StatusOK {
		return nil, status.Errorf(codes.Internal, response.ERR_PUT_FAILED)
	}

	result = &pb.PutConversationResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCConversation) PutLabelToConversation(ctx context.Context, request *pb.PutLabelToConversationRequest) (result *pb.PutLabelToConversationResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	var payload model.ConversationLabelRequest
	if err = util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err = payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	labelId, err := service.PutLabelToConversation(ctx, user, request.GetLabelType(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PutLabelToConversationResponse{
		Code:    "OK",
		Message: "ok",
		Id:      labelId,
	}
	return
}

func (g *GRPCConversation) UpdateStatusConversation(ctx context.Context, request *pb.PutConversationStatusRequest) (*pb.PutConversationStatusResponse, error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	payload := model.ConversationStatusRequest{
		AppId:          request.GetAppId(),
		ConversationId: request.GetConversationId(),
		Status:         request.GetStatus(),
	}
	if err := payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err := service.ConversationService.UpdateStatusConversation(ctx, user, request.GetAppId(), request.GetConversationId(), user.UserId, request.GetStatus())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result := &pb.PutConversationStatusResponse{
		Code:    "OK",
		Message: "ok",
	}
	return result, nil
}

func (g *GRPCConversation) UpdaterUserPreferenceConversation(ctx context.Context, request *pb.UpdaterUserPreferenceConversationRequest) (result *pb.UpdaterUserPreferenceConversationResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	payload := model.ConversationPreferenceRequest{
		AppId:           request.GetAppId(),
		OaId:            request.GetOaId(),
		ConversationId:  request.GetConversationId(),
		PreferenceValue: request.GetPreferenceValue(),
		PreferenceType:  request.GetPreferenceType(),
	}
	if err := payload.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = service.ConversationService.UpdateUserPreferenceConversation(ctx, user, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.UpdaterUserPreferenceConversationResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}
