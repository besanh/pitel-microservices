package repository

import (
	"context"
	"database/sql"

	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
)

type (
	IExternalPluginConnect interface {
		IRepo[model.ExternalPluginConnect]
		GetExternalPluginByType(ctx context.Context, db sqlclient.ISqlClientConn, pluginType string) (result *model.ExternalPluginConnect, err error)
	}
	ExternalPluginConnect struct {
		Repo[model.ExternalPluginConnect]
	}
)

var ExternalPluginConnectRepo IExternalPluginConnect

func NewExternalPluginConnect() IExternalPluginConnect {
	return &ExternalPluginConnect{}
}

func (repo *ExternalPluginConnect) GetExternalPluginByType(ctx context.Context, db sqlclient.ISqlClientConn, pluginType string) (result *model.ExternalPluginConnect, err error) {
	result = new(model.ExternalPluginConnect)
	query := db.GetDB().NewSelect().Model(result).
		Where("plugin_type = ?", pluginType)
	err = query.Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return result, nil
}
