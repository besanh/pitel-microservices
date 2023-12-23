package grpc

import (
	"context"
	"strconv"

	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/common/util"
	pb "github.com/tel4vn/fins-microservices/gen/proto/incom"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCIncom struct {
	pb.UnimplementedIncomServiceServer
}

func NewGRPCIncom() *GRPCIncom {
	return &GRPCIncom{}
}

func (g *GRPCIncom) IncomWebhook(ctx context.Context, request *pb.IncomBodyRequest) (result *pb.IncomBodyResponse, err error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, response.ERR_TOKEN_IS_INVALID)
	}
	authUser := model.AuthUser{
		TenantId:           metadata.Get("tenant_id")[0],
		BusinessUnitId:     metadata.Get("business_unit_id")[0],
		UserId:             metadata.Get("user_id")[0],
		Username:           metadata.Get("username")[0],
		Services:           metadata.Get("services"),
		DatabaseName:       metadata.Get("database_name")[0],
		DatabaseHost:       metadata.Get("database_host")[0],
		DatabasePort:       util.ParseInt(metadata.Get("database_port")[0]),
		DatabaseUser:       metadata.Get("database_user")[0],
		DatabasePassword:   metadata.Get("database_password")[0],
		DatabaseEsHost:     metadata.Get("database_es_host")[0],
		DatabaseEsUser:     metadata.Get("database_es_user")[0],
		DatabaseEsPassword: metadata.Get("database_es_password")[0],
		DatabaseEsIndex:    metadata.Get("database_es_index")[0],
	}
	jsonBody := make(map[string]any)
	if err := util.ParseAnyToAny(request, &jsonBody); err != nil {
		result = &pb.IncomBodyResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}
	var routingConfig string
	routingConfigs := metadata.Get("x_signature_incom")
	if len(routingConfigs) > 0 {
		routingConfig = routingConfigs[0]
	}
	idOmniMess, _ := jsonBody["IdOmniMess"].(string)
	status, _ := jsonBody["Status"].(string)
	channel, _ := jsonBody["Channel"].(string)
	errorCode, _ := jsonBody["ErrorCode"].(string)
	quantityTmp, _ := jsonBody["MtCount"].(float64)
	quantity := int(quantityTmp)
	telcoIdTmp, _ := jsonBody["TelcoId"].(float64)
	telcoId := int(telcoIdTmp)
	isChargedZnsTmp, _ := jsonBody["IsCharged"].(string)
	isChargedZns, _ := strconv.ParseBool(isChargedZnsTmp)

	incomData := model.WebhookIncom{
		IdOmniMess:   idOmniMess,
		Status:       status,
		Channel:      channel,
		ErrorCode:    errorCode,
		Quantity:     quantity,
		TelcoId:      telcoId,
		IsChargedZns: isChargedZns,
	}

	err = service.NewIncom().WebhookIncom(ctx, routingConfig, &authUser, incomData)
	if err != nil {
		result = &pb.IncomBodyResponse{
			Code:    response.MAP_ERR_RESPONSE[response.ERR_DATA_INVALID].Code,
			Message: err.Error(),
		}
		return result, err
	}

	result = &pb.IncomBodyResponse{
		Code:    "OK",
		Message: "success",
	}

	return result, nil
}
