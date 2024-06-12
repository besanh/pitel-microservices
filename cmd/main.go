package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/redis"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/internal/storage"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/server"
	"github.com/tel4vn/fins-microservices/service"
)

var cmdMain = &cobra.Command{
	Use:     "chat-service",
	Short:   "start service",
	Example: "./app chat-service",
	Run: func(cmd *cobra.Command, args []string) {
		RunMainService()
	},
}

func RunMainService() {
	log.InitLogger(config.LogLevel, config.LogFile)

	// init
	initConfig()

	cache.RCache = cache.NewRedisCache(redis.Redis.GetClient())
	defer cache.RCache.Close()

	cache.MCache = cache.NewMemCache()
	defer cache.MCache.Close()

	// Init Repositories
	repository.InitRepositories()
	repository.InitRepositoriesES()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	repository.InitTables(ctx, repository.DBConn)
	repository.InitColumn(ctx, repository.DBConn)
	repository.InitRows(ctx, repository.DBConn)

	// Init services
	service.MapDBConn = make(map[string]sqlclient.ISqlClientConn, 0)
	service.ES_INDEX = env.GetStringENV("ES_INDEX", "pitel_bss_chat")
	service.ES_INDEX_CONVERSATION = env.GetStringENV("ES_INDEX_CONVERSATION", "pitel_bss_conversation")
	service.OTT_URL = env.GetStringENV("OTT_DOMAIN", "")
	service.OTT_VERSION = env.GetStringENV("OTT_VERSION", "v1")
	service.API_SHARE_INFO_HOST = env.GetStringENV("API_SHARE_INFO_HOST", "https://api.dev.fins.vn")
	service.API_CRM = env.GetStringENV("API_CRM", "")
	service.ENABLE_PUBLISH_ADMIN = env.GetBoolENV("ENABLE_PUBLISH_ADMIN", false)
	service.AAA_HOST = env.GetStringENV("AAA_HOST", "https://aaa.dev.fins.vn")
	service.InitServices()

	// Init storage
	storage.InitStorage()

	// Store to service
	service.S3_ENDPOINT = env.GetStringENV("STORAGE_ENDPOINT", "")
	service.S3_BUCKET_NAME = env.GetStringENV("STORAGE_BUCKET_NAME", "")
	service.S3_ACCESS_KEY = env.GetStringENV("STORAGE_ACCESS_KEY", "")
	service.S3_SECRET_KEY = env.GetStringENV("STORAGE_SECRET_KEY", "")

	// Zalo
	service.ZALO_SHARE_INFO_SUBTITLE = env.GetStringENV("ZALO_SHARE_INFO_SUBTITLE", "")

	// Run gRPC server
	server.NewGRPCServer(config.gRPCPort)
}
