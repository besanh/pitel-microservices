package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/tel4vn/fins-microservices/common/constants"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/service/common"
)

type (
	IInboxMarketingFpt interface{}
	InboxMarketingFpt  struct{}
)

func NewInboxMarketingFpt() IInboxMarketingFpt {
	return &InboxMarketingFpt{}
}

func HandleMainInboxMarketingFpt(ctx context.Context, authUser *model.AuthUser, inboxMarketingBasic model.InboxMarketingBasic, routingConfig model.RoutingConfig, inboxMarketing model.InboxMarketingLogInfo, inboxMarketingRequest model.InboxMarketingRequest, fpt model.FptRequireRequest) (model.ResponseInboxMarketing, error) {
	res := model.ResponseInboxMarketing{
		Id: inboxMarketingBasic.DocId,
	}
	dataUpdate := map[string]any{}

	_, result, err := common.HandleDeliveryMessageFpt(ctx, inboxMarketingBasic.DocId, routingConfig, inboxMarketingRequest, fpt)
	if err != nil {
		return res, err
	}
	// Id      string `json:"id"`
	// Status  string `json:"status"`
	// Message string `json:"message"`
	// Code    int    `json:"code"`
	// Find in ES to avoid 404 not found
	dataExist, err := repository.InboxMarketingESRepo.GetDocById(ctx, inboxMarketingBasic.TenantId, authUser.DatabaseEsIndex, inboxMarketingBasic.Id)
	if err != nil {
		return res, err
	} else if len(dataExist.Id) < 1 {
		return res, errors.New("log is not exist")
	}
	// MessageId string `json:"MessageId"`
	// Phone     string `json:"Phone"`
	// BrandName string `json:"BrandName"`
	// Message   string `json:"Message"`
	// PartnerId string `json:"PartnerId"`
	// Telco     string `json:"Telco"`
	var telcoId int
	telcoId, _ = strconv.Atoi(constants.MAP_NETWORK_FPT[result.Telco])
	inboxMarketing.TelcoId = telcoId
	// log
	auditLogModel := model.LogInboxMarketing{
		TenantId:          authUser.TenantId,
		BusinessUnitId:    authUser.BusinessUnitId,
		UserId:            authUser.UserId,
		Username:          authUser.Username,
		Services:          authUser.Services,
		Id:                result.PartnerId,
		RoutingConfigUuid: routingConfig.Id,
		ExternalMessageId: result.MessageId,
		Status:            "",
		Quantity:          0,
		TelcoId:           0,
		IsChargedZns:      false,
		IsCheck:           false,
		Code:              1,
		CountAction:       2,
		UpdatedBy:         inboxMarketingBasic.UpdatedBy,
	}
	auditLog, err := common.HandleAuditLogInboxMarketing(auditLogModel)
	if err != nil {
		return res, err
	}
	inboxMarketing.Log = append(inboxMarketing.Log, auditLog)
	inboxMarketing.UpdatedAt = time.Now()

	tmpBytesUpdate, err := json.Marshal(inboxMarketing)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(tmpBytesUpdate, &dataUpdate); err != nil {
		return res, err
	}

	if err := repository.ESRepo.UpdateDocById(ctx, inboxMarketingBasic.Index, inboxMarketingBasic.DocId, dataUpdate); err != nil {
		return res, err
	}

	res.Status = result.Message

	return res, nil
}
