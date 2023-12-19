package grpc

import (
	"context"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/inbox_marketing_incom"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCInboxMarketingIncom struct {
	pb.UnimplementedIncomServiceServer
}

func NewGRPCInboxMarketingIncom() *GRPCInboxMarketingIncom {
	return &GRPCInboxMarketingIncom{}
}

func (g *GRPCInboxMarketingIncom) HandleWebhook(ctx context.Context, request *pb.IncomBodyRequest) (result *pb.IncomResponse, err error) {
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	log.Info(metaData)
	jsonBody := make(map[string]any)
	if err := util.ParseAnyToAny(request, &jsonBody); err != nil {
		result = &pb.IncomResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	log.Info(jsonBody)

	return nil, nil
}
