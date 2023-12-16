package common

import (
	"context"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/repository"
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
