package service

import (
	"context"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/repository"
)

func CheckBalanceExist(ctx context.Context, dbConn sqlclient.ISqlClientConn, id string) (bool, error) {
	data, err := repository.BalanceConfigRepo.GetById(ctx, dbConn, id)
	if err != nil {
		return false, err
	} else if len(data.Id) < 1 {
		return false, nil
	}
	return true, nil
}
