package grpc

import (
	"context"
	"database/sql"
	"errors"
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
	"google.golang.org/protobuf/types/known/timestamppb"
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

	total, data, respScrollId, err := service.ConversationService.GetConversationsWithScrollAPI(ctx, user, filter, limit, request.GetScrollId())
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ConversationCustomView, 0)
	for _, item := range data {
		var tmp pb.ConversationCustomView
		if err = util.ParseAnyToAny(item, &tmp); err != nil {
			log.Error(err)
			result = &pb.GetConversationsWithScrollAPIResponse{
				Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
				Message: err.Error(),
			}
			return result, nil
		}
		resultData = append(resultData, &tmp)
	}
	result = &pb.GetConversationsWithScrollAPIResponse{
		Code:     "OK",
		Message:  "ok",
		Data:     resultData,
		Total:    int32(total),
		Limit:    int32(limit),
		ScrollId: respScrollId,
	}
	return
}

func (g *GRPCConversation) GetConversationsByManager(ctx context.Context, request *pb.GetConversationsByManagerRequest) (result *pb.GetConversationsByManagerResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		err = status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
		return
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

	total, data, err := service.ConversationService.GetConversationsByHighLevel(ctx, user, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return
	}
	resultData := make([]*pb.ConversationView, 0)
	if data != nil {
		for _, item := range *data {
			tmp, errTmp := convertConversationViewToPbConversationView(&item)
			if errTmp != nil {
				err = errTmp
				log.Error(err)
				result = &pb.GetConversationsByManagerResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return
			}
			resultData = append(resultData, tmp)
		}
	}
	result = &pb.GetConversationsByManagerResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
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

	total, data, respScrollId, err := service.ConversationService.GetConversationsByHighLevelWithScrollAPI(ctx, user, filter, limit, request.GetScrollId())
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ConversationView, 0)
	for _, item := range data {
		tmp, errTmp := convertConversationViewToPbConversationView(item)
		if errTmp != nil {
			err = errTmp
			log.Error(err)
			result = &pb.GetConversationsByManagerWithScrollAPIResponse{
				Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
				Message: err.Error(),
			}
			return
		}
		resultData = append(resultData, tmp)
	}
	result = &pb.GetConversationsByManagerWithScrollAPIResponse{
		Code:     "OK",
		Message:  "ok",
		Data:     resultData,
		Total:    int32(total),
		Limit:    int32(limit),
		ScrollId: respScrollId,
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
	data, err := service.ConversationService.GetConversationById(ctx, user, request.GetAppId(), conversationId)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp, err := convertConversationToPbConversation(data)
	if err != nil {
		log.Error(err)
		result = &pb.GetConversationByIdResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
			Message: err.Error(),
		}
		return
	}
	result = &pb.GetConversationByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data:    tmp,
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
		err = status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
		return
	}
	var payload model.ConversationLabelRequest
	if err = util.ParseAnyToAny(request, &payload); err != nil {
		log.Error(err)
		return
	}

	if err = payload.Validate(); err != nil {
		return
	}

	labelId, err := service.PutLabelToConversation(ctx, user, request.GetLabelType(), payload)
	if err != nil {
		return
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

func (g *GRPCConversation) InsertNoteInConversation(ctx context.Context, request *pb.PostNoteInConversationRequest) (result *pb.PostNoteInConversationResponse, err error) {
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

	id, err := service.ConversationService.InsertNoteInConversation(ctx, user, &payload)
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

func (g *GRPCConversation) UpdateNoteInConversationById(ctx context.Context, request *pb.PutNoteInConversationRequest) (result *pb.PutNoteInConversationResponse, err error) {
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

	err = service.ConversationService.UpdateNoteInConversationById(ctx, user, request.GetNoteId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.PutNoteInConversationResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCConversation) DeleteNoteInConversationById(ctx context.Context, request *pb.DeleteNoteInConversationRequest) (result *pb.DeleteNoteInConversationResponse, err error) {
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
	if err = payload.ValidateDelete(); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = service.ConversationService.DeleteNoteInConversationById(ctx, user, request.GetNoteId(), payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	result = &pb.DeleteNoteInConversationResponse{
		Code:    "OK",
		Message: "ok",
	}
	return
}

func (g *GRPCConversation) GetConversationNotesList(ctx context.Context, request *pb.GetConversationNotesListRequest) (result *pb.GetConversationNotesListResponse, err error) {
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

	total, data, err := service.ConversationService.GetConversationNotesList(ctx, user, filter, limit, offset)
	if err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.NotesList, 0)
	for _, item := range data {
		tmp := pb.NotesList{
			Id:        item.Id,
			Content:   item.Content,
			CreatedAt: timestamppb.New(item.CreatedAt),
			UpdatedAt: timestamppb.New(item.UpdatedAt),
		}

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
