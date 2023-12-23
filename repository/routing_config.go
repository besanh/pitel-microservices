package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/model"
)

type (
	IRoutingConfig interface {
		GetRoutingConfigById(ctx context.Context, routingConfigUuid string) (*model.RoutingConfig, error)
	}
	RoutingConfig struct {
	}
)

var RoutingConfigRepo IRoutingConfig

func NewRoutingConfig() IRoutingConfig {
	return &RoutingConfig{}
}

func (repo *RoutingConfig) GetRoutingConfigById(ctx context.Context, routingConfigUuid string) (result *model.RoutingConfig, err error) {
	result = new(model.RoutingConfig)
	query := DBConn.
		GetDB().
		NewSelect().
		Model(result).
		Where("id = ?", routingConfigUuid)
	if err = query.Scan(ctx); err == sql.ErrNoRows {
		return nil, nil
	}
	return
}
