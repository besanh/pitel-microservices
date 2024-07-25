package grpc

import (
	"errors"
	"net/http"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/conversation"
	"github.com/tel4vn/fins-microservices/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func parseResponseDataOfGetConversationsWithScrollAPI(code int, responseData any, limit int) (result *pb.GetConversationsWithScrollAPIResponse, err error) {
	if code != http.StatusOK {
		err = errors.New(response.ERR_GET_FAILED)
		return
	}
	var paginationData map[string]any
	if err = util.ParseAnyToAny(responseData, &paginationData); err != nil {
		return
	}

	total, ok2 := paginationData["total"].(float64)
	if !ok2 {
		err = errors.New("not found total")
		return
	}

	var parsedData map[string]any
	if err = util.ParseAnyToAny(paginationData["data"], &parsedData); err != nil {
		return
	}
	respScrollId, ok3 := parsedData["scroll_id"].(string)
	if !ok3 {
		err = errors.New("scroll_id not found in parsedData")
		return
	}
	var data []model.ConversationCustomView
	if err = util.ParseAnyToAny(parsedData["conversations"], &data); err != nil {
		return
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

func parseResponseDataOfGetConversationsByManager(code int, responseData any, limit, offset int) (result *pb.GetConversationsByManagerResponse, err error) {
	if code != http.StatusOK {
		err = errors.New(response.ERR_GET_FAILED)
		return
	}
	var paginationData map[string]any
	if err = util.ParseAnyToAny(responseData, &paginationData); err != nil {
		return
	}
	total, ok2 := paginationData["total"].(float64)
	if !ok2 {
		err = errors.New("not found total")
		return
	}
	var data *[]model.ConversationView
	if err = util.ParseAnyToAny(paginationData["data"], &data); err != nil {
		return
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

func parseResponseDataOfGetConversationsByManagerWithScrollAPI(code int, responseData any, limit int) (result *pb.GetConversationsByManagerWithScrollAPIResponse, err error) {
	if code != http.StatusOK {
		err = errors.New(response.ERR_GET_FAILED)
		return
	}
	var paginationData map[string]any
	if err = util.ParseAnyToAny(responseData, &paginationData); err != nil {
		return
	}
	total, ok2 := paginationData["total"].(float64)
	if !ok2 {
		err = errors.New("not found total")
		return
	}
	var parsedData map[string]any
	if err = util.ParseAnyToAny(paginationData["data"], &parsedData); err != nil {
		return
	}
	respScrollId, ok3 := parsedData["scroll_id"].(string)
	if !ok3 {
		err = errors.New("scroll_id not found in parsedData")
		return
	}
	var data []*model.ConversationView
	if err = util.ParseAnyToAny(parsedData["conversations"], &data); err != nil {
		return
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

func parseResponseDataOfGetConversationById(code int, responseData any) (result *pb.GetConversationByIdResponse, err error) {
	if code != http.StatusOK {
		err = errors.New(response.ERR_GET_FAILED)
		return
	}
	var data model.Conversation
	if err = util.ParseAnyToAny(responseData, &data); err != nil {
		return
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
			return
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
