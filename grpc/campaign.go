package grpc

import (
	"context"
	"fmt"
	"log"

	validator "github.com/bufbuild/protovalidate-go"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/example"
	"github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCExample struct {
	pb.UnimplementedExampleServiceServer
}

func NewGRPCExample() *GRPCExample {
	return &GRPCExample{}
}

func (s *GRPCExample) GetExamples(ctx context.Context, req *pb.GetExamplesRequest) (result *pb.GetExamplesResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	limit := util.ParseLimit(req.GetLimit())
	offset := util.ParseOffset(req.GetOffset())
	params := make([]model.Param, 0)
	params = append(params, model.Param{
		Key:      "example_name",
		Value:    fmt.Sprintf("%%%s%%", req.GetKeyword()),
		Operator: "LIKE",
	})
	var data []*structpb.Struct
	examples, total, err := service.ExampleService.GetExamples(ctx, authUser, params, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	} else if len(*examples) > 0 {
		for _, campaign := range *examples {
			element := map[string]any{
				"example_name": campaign.ExampleName,
				"id":           campaign.Id,
				"created_at":   campaign.CreatedAt,
				"updated_at":   campaign.UpdatedAt,
			}
			tmp, _ := util.ToStructPb(element)
			data = append(data, tmp)
		}
	}
	result = &pb.GetExamplesResponse{
		Code:    "OK",
		Message: "ok",
		Data:    data,
		Total:   int32(total),
	}
	return
}

func (c *GRPCExample) GetExampleById(ctx context.Context, req *pb.GetExampleByIdRequest) (result *pb.GetExampleByIdResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	example, err := service.ExampleService.GetExampleById(ctx, authUser, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	result = &pb.GetExampleByIdResponse{
		Code:    "OK",
		Message: "ok",
		Data: &pb.ExampleResponseData{
			Id:          example.Id,
			ExampleName: example.ExampleName,
			CreatedAt:   timestamppb.New(example.CreatedAt),
			UpdatedAt:   timestamppb.New(example.UpdatedAt),
		},
	}
	return
}

func (c *GRPCExample) PostExample(ctx context.Context, req *pb.PostExampleRequest) (result *pb.PostExampleResponse, err error) {
	validator, err := validator.New()
	if err != nil {
		log.Fatal(err)
	}
	err = validator.Validate(req.GetData())
	if err != nil {
		result = &pb.PostExampleResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_EXAMPLE_INVALID].Code,
			Message: err.Error(),
		}
		if isStr, r := response.HandleValidatorPBError(err); isStr {
			result.ErrorCode = &pb.PostExampleResponse_Error{Error: r.(string)}
		} else {
			var res *structpb.Struct
			res, err = util.ToStructPb(r)
			if err != nil {
				return nil, status.Error(codes.Unavailable, err.Error())
			}
			result.ErrorCode = &pb.PostExampleResponse_Errors{Errors: res}
		}
		return result, nil
	}
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	example := model.Example{
		ExampleName: req.GetData().GetExampleName(),
	}
	if err := service.ExampleService.PostExample(ctx, authUser, example); err != nil {
		return nil, status.Error(codes.Unavailable, err.Error())
	}
	var res *structpb.Struct
	camp := map[string]any{
		"id": example.Id,
	}
	res, _ = util.ToStructPb(camp)
	result = &pb.PostExampleResponse{
		Code:    "OK",
		Message: "ok",
		Data:    res,
	}
	return
}

func (c *GRPCExample) PutExample(ctx context.Context, req *pb.PutExampleRequest) (result *pb.PutExampleResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	_ = authUser
	var res *structpb.Struct
	camp := map[string]any{
		"id": req.GetId(),
	}
	res, _ = util.ToStructPb(camp)
	result = &pb.PutExampleResponse{
		Code:    "OK",
		Message: "ok",
		Data:    res,
	}
	return
}

func (c *GRPCExample) DeleteExample(ctx context.Context, req *pb.DeleteExampleRequest) (result *pb.DeleteExampleResponse, err error) {
	authUser, ok := auth.GetUserFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	_ = authUser
	var res *structpb.Struct
	res, _ = util.ToStructPb(map[string]any{
		"id": req.GetId(),
	})
	result = &pb.DeleteExampleResponse{
		Code:    "OK",
		Message: "ok",
		Data:    res,
	}
	return
}
