package common

import (
	"context"
	"time"

	cacheUtil "github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/model"
	"github.com/tel4vn/fins-microservices/repository"
)

const (
	INFO_EXTERNAL_PLUGIN_CONNECT   = "INFO_EXTERNAL_PLUGIN_CONNECT"
	EXPIRE_EXTERNAL_PLUGIN_CONNECT = 1 * time.Minute
)

func GetExternalPluginConnectFromCache(ctx context.Context, dbCon sqlclient.ISqlClientConn, externalPluginConnectType string) (*model.ExternalPluginConnect, error) {
	externalPluginConnectCache, err := cacheUtil.MCache.Get(INFO_EXTERNAL_PLUGIN_CONNECT + "_" + externalPluginConnectType)
	if err != nil {
		return nil, err
	} else if externalPluginConnectCache != nil {
		externalPluginConnect := externalPluginConnectCache.(*model.ExternalPluginConnect)
		return externalPluginConnect, nil
	} else {
		externalPluginConnect, err := repository.ExternalPluginConnectRepo.GetExternalPluginByType(ctx, dbCon, externalPluginConnectType)
		if err != nil {
			return nil, err
		}
		if err := cacheUtil.MCache.SetTTL(INFO_EXTERNAL_PLUGIN_CONNECT+"_"+externalPluginConnectType, externalPluginConnect, EXPIRE_EXTERNAL_PLUGIN_CONNECT); err != nil {
			return nil, err
		}
		return externalPluginConnect, nil
	}
}
