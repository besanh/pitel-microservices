package grpc

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/inbox_marketing"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type GRPCInboxMarketing struct {
	pb.UnsafeInboxMarketingServiceServer
}

func NewGRPCInboxMarketing() *GRPCInboxMarketing {
	return &GRPCInboxMarketing{}
}

func (g *GRPCInboxMarketing) SendInboxMarketing(ctx context.Context, request *pb.InboxMarketingBodyRequest) (result *pb.InboxMarketingResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	var payload model.InboxMarketingRequest
	if err := util.ParseAnyToAny(request, &payload); err != nil {
		result = &pb.InboxMarketingResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Status:  "failed",
			Message: err.Error(),
		}
		return result, err
	}

	if err := payload.Validate(); err != nil {
		log.Error(err)
		result = &pb.InboxMarketingResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Status:  "failed",
			Message: err.Error(),
		}
		return result, nil
	}

	res, err := service.NewInboxMarketing().SendInboxMarketing(ctx, authUser, payload)
	if err != nil {
		result = &pb.InboxMarketingResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Status:  res.Status,
			Message: res.Message,
		}
		return result, nil
	}

	result = &pb.InboxMarketingResponse{
		Code:    res.Code,
		Status:  res.Status,
		Message: res.Message,
	}

	return result, nil
}

func (g *GRPCInboxMarketing) ReportInboxMarketing(ctx context.Context, request *pb.InboxMarketingRequest) (result *pb.InboxMarketingStructResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}

	limit := util.ParseLimit(request.GetLimit())
	offset := util.ParseOffset(request.GetOffset())

	isChargedZnsQuery, _ := strconv.ParseBool(request.GetIsChargedZns())
	isChargedZns := sql.NullBool{}
	if len(request.GetIsChargedZns()) > 0 {
		isChargedZns.Bool = isChargedZnsQuery
		isChargedZns.Valid = true
	}
	filter := model.InboxMarketingFilter{
		StartTime:         request.GetStartTime(),
		EndTime:           request.GetEndTime(),
		TenantId:          request.GetTenantId(),
		BusinessUnitId:    request.GetBusinessUnitId(),
		UserId:            request.GetUserId(),
		Username:          request.GetUsername(),
		Services:          request.GetSevices(),
		RoutingConfigUuid: request.GetRoutingConfigUuid(),
		Plugin:            request.GetPlugin(),
		PhoneNumber:       request.GetPhoneNumber(),
		Message:           request.GetMessage(),
		TemplateUuid:      request.GetTemplateUuid(),
		TemplateCode:      request.GetTemplateCode(),
		Channel:           request.GetChannel(),
		Status:            request.GetStatus(),
		ErrorCode:         request.GetErrorCode(),
		Quantity:          request.GetQuantity(),
		TelcoId:           request.GetTelcoId(),
		RouteRule:         request.GetRouteRule(),
		ServiceTypeId:     request.GetServiceTypeId(),
		SendTime:          request.GetSendTime(),
		Ext:               request.GetExt(),
		IsChargedZns:      isChargedZns,
		Code:              request.GetCode(),
		CountAction:       request.GetCountAction(),
		CampaignUuid:      request.GetCampaignUuid(),
	}

	total, res, err := service.NewInboxMarketing().GetReportInboxMarketing(ctx, authUser, filter, limit, offset)
	if err != nil {
		result = &pb.InboxMarketingStructResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, nil
	}
	var data []*structpb.Struct
	if len(res) > 0 {
		for _, val := range res {
			var item model.InboxMarketingLogReport
			if err = util.ParseAnyToAny(val, &item); err != nil {
				result = &pb.InboxMarketingStructResponse{
					Code:    response.MAP_ERR_RESPONSE[response.ERR_CAMPAIGN_INVALID].Code,
					Message: err.Error(),
				}
				return result, nil
			}

			tmp, _ := util.ToStructPb(item)
			data = append(data, tmp)
		}
	}

	result = &pb.InboxMarketingStructResponse{
		Code:    "OK",
		Message: "SUCCESS",
		Data:    data,
		Total:   int32(total),
	}

	return
}
