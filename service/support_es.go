package service

import (
	"context"
	"encoding/json"

	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/repository"
)

func InsertES(ctx context.Context, appId, index, docId string, data any) error {
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
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, index, appId); err != nil {
		log.Error(err)
		return err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, index, appId); err != nil {
			log.Error(err)
			return err
		}
	}

	err = repository.ESRepo.InsertLog(ctx, appId, ES_INDEX, docId, esDoc)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
