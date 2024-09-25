package grpc

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/chat_report"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCChatReport struct{}

func NewGRPCChatReport() pb.ChatReportServiceServer {
	return &GRPCChatReport{}
}

func (g *GRPCChatReport) GetWorkPerformanceReportByUser(ctx context.Context, request *pb.GetWorkPerformanceReportByUserRequest) (result *pb.GetWorkPerformanceReportByUserResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	filter := model.MessageFilter{
		StartTime:   request.GetStartTime(),
		EndTime:     request.GetEndTime(),
		SupporterId: request.GetUserId(),
		UnitUuid:    request.GetUnitUuid(),
	}

	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	total, data, err := service.ChatReportService.GetChatWorkReports(ctx, user, filter, limit, offset, token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.WorkReport, 0)
	if data != nil {
		resultData = convertChatWorkReportToPbChatWorkReport(data)
	}
	result = &pb.GetWorkPerformanceReportByUserResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return
}

func (g *GRPCChatReport) GetMultichannelPerformanceReport(ctx context.Context, request *pb.GetMultichannelPerformanceReportRequest) (result *pb.GetMultichannelPerformanceReportResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	filter := model.ConversationFilter{
		StartTime: request.GetStartTime(),
		EndTime:   request.GetEndTime(),
	}
	if len(request.GetChannel()) > 0 {
		filter.ConversationType = []string{request.GetChannel()}
	}
	if len(request.GetOaId()) > 0 {
		filter.OaIds = []string{request.GetOaId()}
	}

	limit, offset := util.ParseLimit(request.GetLimit()), util.ParseOffset(request.GetOffset())

	total, data, err := service.ChatReportService.GetChatGeneralReports(ctx, user, filter, limit, offset, token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resultData := make([]*pb.GeneralReport, 0)
	if data != nil {
		resultData = convertChatGeneralReportToPbChatGeneralReport(data)
	}
	result = &pb.GetMultichannelPerformanceReportResponse{
		Code:    "OK",
		Message: "ok",
		Data:    resultData,
		Total:   int32(total),
		Limit:   int32(limit),
		Offset:  int32(offset),
	}
	return
}

func (g *GRPCChatReport) ExportWorkPerformanceReportByUser(ctx context.Context, request *pb.ExportWorkPerformanceReportByUserRequest) (result *pb.ExportWorkPerformanceReportByUserResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	filter := model.MessageFilter{
		StartTime:   request.GetStartTime(),
		EndTime:     request.GetEndTime(),
		SupporterId: request.GetUserId(),
	}
	if !slices.Contains([]string{"xlsx", "csv"}, request.GetFileType()) {
		return nil, status.Errorf(codes.InvalidArgument, "file type not supported")
	}

	data, err := service.ChatReportService.ExportWorkReports(ctx, user, filter, request.GetFileType(), token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	result = &pb.ExportWorkPerformanceReportByUserResponse{
		Code:     "OK",
		Message:  "ok",
		FileName: data,
	}
	return
}

func (g *GRPCChatReport) ExportMultichannelPerformanceReport(ctx context.Context, request *pb.ExportMultichannelPerformanceReportRequest) (result *pb.ExportMultichannelPerformanceReportResponse, err error) {
	user, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}

	filter := model.ConversationFilter{
		StartTime: request.GetStartTime(),
		EndTime:   request.GetEndTime(),
	}
	if len(request.GetChannel()) > 0 {
		filter.ConversationType = []string{request.GetChannel()}
	}
	if len(request.GetOaId()) > 0 {
		filter.OaIds = []string{request.GetOaId()}
	}

	if !slices.Contains([]string{"xlsx", "csv"}, request.GetFileType()) {
		return nil, status.Errorf(codes.InvalidArgument, "file type not supported")
	}

	data, err := service.ChatReportService.ExportGeneralReports(ctx, user, filter, request.GetFileType(), token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	result = &pb.ExportMultichannelPerformanceReportResponse{
		Code:     "OK",
		Message:  "ok",
		FileName: data,
	}
	return
}
