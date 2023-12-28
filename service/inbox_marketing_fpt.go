package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

type (
	IInboxMarketingFpt interface{}
	InboxMarketingFpt  struct{}
)

func NewInboxMarketingFpt() IInboxMarketingFpt {
	return &InboxMarketingFpt{}
}

func HandleMainInboxMarketingFpt(ctx context.Context, authUser *model.AuthUser, inboxMarketingBasic model.InboxMarketingBasic, routingConfig model.RoutingConfig, inboxMarketing model.InboxMarketingLogInfo, inboxMarketingRequest model.InboxMarketingRequest, fpt model.FptRequireRequest) (model.ResponseInboxMarketing, error) {
	dataUpdate := map[string]any{}
	statusCode, result, resultFpt, err := HandleDeliveryMessageFpt(ctx, inboxMarketingBasic.DocId, routingConfig, inboxMarketingRequest, fpt)
	if err != nil {
		log.Error(err)
		return *result, err
	} else if statusCode != 200 {
		return *result, errors.New(result.Message)
	}

	// Find in ES to avoid 404 not found
	dataExist, err := repository.InboxMarketingESRepo.GetDocById(ctx, inboxMarketingBasic.TenantId, ES_INDEX, inboxMarketingBasic.DocId)
	if err != nil {
		return *result, err
	} else if len(dataExist.Id) < 1 {
		return *result, errors.New("log is not exist")
	}

	var telcoId int
	telcoId, _ = strconv.Atoi(constants.MAP_NETWORK_FPT[resultFpt.Telco])
	inboxMarketing.TelcoId = telcoId
	inboxMarketing.ExternalMessageId = resultFpt.MessageId
	// log
	auditLogModel := model.LogInboxMarketing{
		TenantId:          authUser.TenantId,
		BusinessUnitId:    authUser.BusinessUnitId,
		UserId:            authUser.UserId,
		Username:          authUser.Username,
		Services:          authUser.Services,
		Id:                inboxMarketing.Id,
		RoutingConfigUuid: routingConfig.Id,
		ExternalMessageId: resultFpt.MessageId,
		Status:            result.Status,
		Quantity:          0,
		TelcoId:           0,
		IsChargedZns:      false,
		IsCheck:           false,
		Code:              1,
		CountAction:       2,
		UpdatedBy:         inboxMarketingBasic.UpdatedBy,
	}
	auditLog, err := HandleAuditLogInboxMarketing(auditLogModel)
	if err != nil {
		return *result, err
	}
	inboxMarketing.Log = append(inboxMarketing.Log, auditLog)
	inboxMarketing.UpdatedAt = time.Now()

	tmpBytesUpdate, err := json.Marshal(inboxMarketing)
	if err != nil {
		return *result, err
	}
	if err := json.Unmarshal(tmpBytesUpdate, &dataUpdate); err != nil {
		return *result, err
	}

	if err := repository.ESRepo.UpdateDocById(ctx, inboxMarketingBasic.Index, inboxMarketingBasic.DocId, dataUpdate); err != nil {
		return *result, err
	}

	return *result, nil
}
