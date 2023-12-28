package service

import (
	"context"
	"time"

	cacheUtil "github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

const (
	INFO_RECIPIENT   = "info_recipient"
	EXPIRE_RECIPIENT = 30 * time.Minute
)

func CheckRecipientExist(ctx context.Context, dbConn sqlclient.ISqlClientConn, id string) (bool, error) {
	data, err := repository.RecipientConfigRepo.GetById(ctx, dbConn, id)
	if err != nil {
		return false, err
	} else if len(data.Id) < 1 {
		return false, nil
	}
	return true, nil
}

// Get route rule from recipient config
func HandleGetRouteRule(ctx context.Context, dbConn sqlclient.ISqlClientConn, recipientId string) (routeRules string, err error) {
	recipientConfigCache := cacheUtil.NewMemCache().Get(INFO_RECIPIENT + "_" + recipientId)
	if recipientConfigCache != nil {
		recipientConfig := recipientConfigCache.(*model.RecipientConfig)
		if recipientConfig.Provider == "incom" {
			routeRules = routingIncomRoutRule(recipientConfig.RecipientType)
		}
		return routeRules, nil
	} else {
		recipientConfig, err := repository.RecipientConfigRepo.GetById(ctx, dbConn, recipientId)
		if err != nil {
			return "", err
		}
		if recipientConfig.Provider == "incom" {
			routeRules = routingIncomRoutRule(recipientConfig.RecipientType)
		}
		cacheUtil.NewMemCache().Set(INFO_RECIPIENT+"_"+recipientId, recipientConfig, EXPIRE_RECIPIENT)
		return routeRules, nil
	}
}

// Note: current rule incom: 1: ZNS , 2: Autocall, 3: SMS
func routingIncomRoutRule(recipientType string) (routeRule string) {
	if recipientType == "zns" {
		routeRule = "1"
	} else if recipientType == "autocall" {
		routeRule = "2"
	} else if recipientType == "sms" {
		routeRule = "3"
	}
	return
}

func removeDuplicateValues(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
