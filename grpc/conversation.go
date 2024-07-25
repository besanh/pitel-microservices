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

	code, responseData := service.ConversationService.GetConversationsWithScrollAPI(ctx, user, filter, limit, request.GetScrollId())
	if code != http.StatusOK {
		return nil, status.Errorf(codes.Internal, response.ERR_GET_FAILED)
	}
	var paginationData map[string]any
	if err = util.ParseAnyToAny(responseData, &paginationData); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	total, ok2 := paginationData["total"].(float64)
	if !ok2 {
		log.Error("not found total")
		return nil, status.Errorf(codes.Internal, "not found total")
	}

	var parsedData map[string]any
	if err = util.ParseAnyToAny(paginationData["data"], &parsedData); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	respScrollId, ok3 := parsedData["scroll_id"].(string)
	if !ok3 {
		log.Error("scroll_id not found in parsedData")
		return nil, status.Errorf(codes.Internal, "scroll_id not found in parsedData")
	}
	var data []model.ConversationCustomView
	if err = util.ParseAnyToAny(parsedData["conversations"], &data); err != nil {
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
	if code != http.StatusOK {
		return nil, status.Errorf(codes.Internal, response.ERR_GET_FAILED)
	}
	var paginationData map[string]any
	if err = util.ParseAnyToAny(responseData, &paginationData); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	total, ok2 := paginationData["total"].(float64)
	if !ok2 {
		log.Error("not found total")
		return nil, status.Errorf(codes.Internal, "not found total")
	}
	var data *[]model.ConversationView
	if err = util.ParseAnyToAny(paginationData["data"], &data); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.ConversationView, 0)
	if data != nil {
		for _, item := range *data {
			tmp := &pb.ConversationView{
				TenantId:               item.TenantId,
				ConversationId:         item.ConversationId,
				ConversationType:       item.ConversationType,
				AppId:                  item.AppId,
				OaId:                   item.OaId,
				OaName:                 item.OaName,
				OaAvatar:               item.OaAvatar,
				ShareInfo:              nil,
				ExternalUserId:         item.ExternalUserId,
				Username:               item.Username,
				Avatar:                 item.Avatar,
				Major:                  item.Major,
				Following:              item.Following,
				Labels:                 make([]*pb.ChatLabel, 0),
				IsDone:                 item.IsDone,
				IsDoneAt:               item.IsDoneAt,
				IsDoneBy:               item.IsDoneBy,
				CreatedAt:              item.CreatedAt,
				UpdatedAt:              item.UpdatedAt,
				TotalUnread:            item.TotalUnRead,
				LatestMessageContent:   item.LatestMessageContent,
				LatestMessageDirection: item.LatestMessageDirection,
				ExternalConversationId: item.ExternalConversationId,
			}
			if item.ShareInfo != nil {
				if err = util.ParseAnyToAny(item.ShareInfo, &tmp.ShareInfo); err != nil {
					log.Error(err)
					result = &pb.GetConversationsByManagerResponse{
						Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
						Message: err.Error(),
					}
					return
				}
			}
			if err = util.ParseStringToAny(string(item.Label), &tmp.Labels); err != nil {
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

	code, responseData := service.ConversationService.GetConversationsByHighLevelWithScrollAPI(ctx, user, filter, limit, request.GetScrollId())
	if code != http.StatusOK {
		return nil, status.Errorf(codes.Internal, response.ERR_GET_FAILED)
	}
	var paginationData map[string]any
	if err = util.ParseAnyToAny(responseData, &paginationData); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	total, ok2 := paginationData["total"].(float64)
	if !ok2 {
		log.Error("not found total")
		return nil, status.Errorf(codes.Internal, "not found total")
	}
	var parsedData map[string]any
	if err = util.ParseAnyToAny(paginationData["data"], &parsedData); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	respScrollId, ok3 := parsedData["scroll_id"].(string)
	if !ok3 {
		log.Error("scroll_id not found in parsedData")
		return nil, status.Errorf(codes.Internal, "scroll_id not found in parsedData")
	}
	var data []*model.ConversationView
	if err = util.ParseAnyToAny(parsedData["conversations"], &data); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resultData := make([]*pb.ConversationView, 0)
	for _, item := range data {
		tmp := &pb.ConversationView{
			TenantId:               item.TenantId,
			ConversationId:         item.ConversationId,
			ConversationType:       item.ConversationType,
			AppId:                  item.AppId,
			OaId:                   item.OaId,
			OaName:                 item.OaName,
			OaAvatar:               item.OaAvatar,
			ShareInfo:              nil,
			ExternalUserId:         item.ExternalUserId,
			Username:               item.Username,
			Avatar:                 item.Avatar,
			Major:                  item.Major,
			Following:              item.Following,
			Labels:                 make([]*pb.ChatLabel, 0),
			IsDone:                 item.IsDone,
			IsDoneAt:               item.IsDoneAt,
			IsDoneBy:               item.IsDoneBy,
			CreatedAt:              item.CreatedAt,
			UpdatedAt:              item.UpdatedAt,
			TotalUnread:            item.TotalUnRead,
			LatestMessageContent:   item.LatestMessageContent,
			LatestMessageDirection: item.LatestMessageDirection,
			ExternalConversationId: item.ExternalConversationId,
		}
		if item.ShareInfo != nil {
			if err = util.ParseAnyToAny(item.ShareInfo, &tmp.ShareInfo); err != nil {
				log.Error(err)
				result = &pb.GetConversationsByManagerWithScrollAPIResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_GET_FAILED].Code,
					Message: err.Error(),
				}
				return
			}
		}
		if err = util.ParseStringToAny(string(item.Label), &tmp.Labels); err != nil {
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
	code, responseData := service.ConversationService.GetConversationById(ctx, user, request.GetAppId(), conversationId)
	if code != http.StatusOK {
		return nil, status.Errorf(codes.Internal, response.ERR_GET_FAILED)
	}
	var data model.Conversation
	if err = util.ParseAnyToAny(responseData, &data); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	tmp := &pb.Conversation{
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
	}
	if data.ShareInfo != nil {
		if err = util.ParseAnyToAny(data.ShareInfo, &tmp.ShareInfo); err != nil {
			log.Error(err)
			return nil, status.Errorf(codes.Internal, err.Error())
		}
	}
	if err = util.ParseStringToAny(string(data.Labels), &tmp.Labels); err != nil {
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

	code, responseData := service.ConversationService.UpdateConversationById(ctx, user, request.GetAppId(), request.GetOaId(), request.GetId(), payload)
	if code != http.StatusOK {
		return nil, status.Errorf(codes.Internal, response.ERR_PUT_FAILED)
	}
	var data map[string]any
	if err = util.ParseAnyToAny(responseData, &data); err != nil {
		log.Error(err)
		return nil, status.Errorf(codes.Internal, err.Error())
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
