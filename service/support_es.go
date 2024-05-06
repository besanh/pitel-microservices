package service

import (
	"context"
	"encoding/json"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/repository"
)

func InsertES(ctx context.Context, tenantId, index, appId, docId string, data any) error {
	tmpBytes, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		return err
	}
	esDoc := map[string]any{}
	if err := json.Unmarshal(tmpBytes, &esDoc); err != nil {
		log.Error(err)
		return err
	}
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, index, tenantId); err != nil {
		log.Error(err)
		return err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, index, tenantId); err != nil {
			log.Error(err)
			return err
		}
	}

	if err = repository.ESRepo.InsertLog(ctx, tenantId, ES_INDEX, appId, docId, esDoc); err != nil {
		log.Error(err)
		return err
	}
	return nil
}
