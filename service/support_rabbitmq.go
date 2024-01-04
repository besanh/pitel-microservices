package service

import (
	"context"
	"encoding/json"

	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

func HandlePushRMQ(ctx context.Context, index, docId string, message model.Message, tmpBytes []byte) error {
	esDoc := make(map[string]any)
	err := json.Unmarshal(tmpBytes, &esDoc)
	if err != nil {
		return err
	}
	if isExisted, err := repository.ESRepo.CheckAliasExist(ctx, ES_INDEX, message.AppId); err != nil {
		return err
	} else if !isExisted {
		if err := repository.ESRepo.CreateAlias(ctx, ES_INDEX, message.AppId); err != nil {
			return err
		}
	}

	_, err = repository.ESRepo.CreateDocRabbitMQ(ctx, ES_INDEX, message.AppId, message.AppId, docId, esDoc)
	if err != nil {
		return err
	}

	return nil
}
