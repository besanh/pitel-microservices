package cmd

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/minio"
	"github.com/tel4vn/fins-microservices/internal/redis"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	"github.com/tel4vn/fins-microservices/internal/storage"
	"github.com/tel4vn/fins-microservices/repository"
	"github.com/tel4vn/fins-microservices/server"
	"github.com/tel4vn/fins-microservices/service"
)

var cmdMain = &cobra.Command{
	Use:     "chat",
	Short:   "start service",
	Example: "./app chat-service",
	Run: func(cmd *cobra.Command, args []string) {
		RunMainService()
	},
}

func RunMainService() {
	log.InitLogger(config.LogLevel, config.LogFile)

	// init cache
	cache.RCache = cache.NewRedisCache(redis.Redis.GetClient())
	cache.MCache = cache.NewMemCache()

	// Init Repositories
	repository.InitRepositories()
	repository.InitRepositoriesES()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Init tables and columns
	repository.InitTables(ctx, repository.DBConn)
	repository.InitColumn(ctx, repository.DBConn)

	// Init services
	service.SECRET_KEY_SUPERADMIN = env.GetStringENV("SECRET_KEY_SUPERADMIN", "RnXyO4178f2gvXV8bbSgVf3ipcO7PR5y6jLATrfvHcmEbVWjgwgm2dl8GE3EPEG7KFqHzOznCNBbe3aiNWykfT32lw0RM8ThRTCD")
	service.MapDBConn = make(map[string]sqlclient.ISqlClientConn, 0)
	service.ES_INDEX = env.GetStringENV("ES_INDEX", "pitel_bss_chat")
	service.ES_INDEX_CONVERSATION = env.GetStringENV("ES_INDEX_CONVERSATION", "pitel_bss_conversation")
	service.OTT_URL = env.GetStringENV("OTT_DOMAIN", "")
	service.OTT_VERSION = env.GetStringENV("OTT_VERSION", "v1")
	service.API_SHARE_INFO_HOST = env.GetStringENV("API_SHARE_INFO_HOST", "")
	service.API_DOC = env.GetStringENV("API_DOC", "")
	service.ENABLE_PUBLISH_ADMIN = env.GetBoolENV("ENABLE_PUBLISH_ADMIN", false)
	service.ENABLE_CHAT_AUTO_SCRIPT_REPLY = env.GetBoolENV("ENABLE_CHAT_AUTO_SCRIPT_REPLY", false)
	service.ENABLE_CHAT_POLICY_SETTINGS = env.GetBoolENV("ENABLE_CHAT_POLICY_SETTINGS", false)
	service.AAA_HOST = env.GetStringENV("AAA_HOST", "https://aaa.dev.fins.vn")
	service.InitServices()

	// Init storage
	storage.InitStorage()

	// Store to service
	minio.MinIOClient = minio.NewClient(minio.Config{
		Endpoint:        env.GetStringENV("STORAGE_ENDPOINT", ""),
		AccessKeyID:     env.GetStringENV("STORAGE_BUCKET_NAME", ""),
		SecretAccessKey: env.GetStringENV("STORAGE_ACCESS_KEY", ""),
		Region:          env.GetStringENV("STORAGE_SECRET_KEY", ""),
		UseSSL:          true,
	})

	initConfigService()

	// Run gRPC server
	log.Debug("run gRPC server")
	server.NewGRPCServer(config.gRPCPort)
}

func initConfigService() {
	service.S3_ENDPOINT = env.GetStringENV("STORAGE_ENDPOINT", "")
	service.S3_BUCKET_NAME = env.GetStringENV("STORAGE_BUCKET_NAME", "")
	service.S3_ACCESS_KEY = env.GetStringENV("STORAGE_ACCESS_KEY", "")
	service.S3_SECRET_KEY = env.GetStringENV("STORAGE_SECRET_KEY", "")

	// Zalo
	service.ZALO_SHARE_INFO_SUBTITLE = env.GetStringENV("ZALO_SHARE_INFO_SUBTITLE", "")
	service.ZALO_POLICY_CHAT_WINDOW = env.GetIntENV("ZALO_POLICY_CHAT_WINDOW", 604800) // 7 days in secs

	// Facebook
	service.FACEBOOK_GRAPH_API_VERSION = env.GetStringENV("FACEBOOK_GRAPH_API_VERSION", "")
	service.FACEBOOK_POLICY_CHAT_WINDOW = env.GetIntENV("FACEBOOK_POLICY_CHAT_WINDOW", 86400) // 1 day in secs

	// DB for cronjob
	service.DB_HOST = env.GetStringENV("DB_HOST", "")
	service.DB_DATABASE = env.GetStringENV("DB_DATABASE", "")
	service.DB_USERNAME = env.GetStringENV("DB_USERNAME", "")
	service.DB_PASSWORD = env.GetStringENV("DB_PASSWORD", "")
	service.DB_PORT = env.GetIntENV("DB_PORT", 0)

	// Smtp
	service.SMTP_SERVER = env.GetStringENV("SMTP_SERVER", "")
	service.SMTP_MAILPORT = env.GetIntENV("SMTP_MAILPORT", 465)
	service.SMTP_USERNAME = env.GetStringENV("SMTP_USERNAME", "")
	service.SMTP_PASSWORD = env.GetStringENV("SMTP_PASSWORD", "")
	service.SMTP_INFORM = env.GetBoolENV("SMTP_INFORM", false)
	service.ENABLE_NOTIFY_EMAIL = env.GetBoolENV("ENABLE_NOTIFY_EMAIL", false)

	if service.ENABLE_NOTIFY_EMAIL {
		log.Info("init scheduler for expire token")
		s1 := gocron.NewScheduler(time.Local)
		s1.SetMaxConcurrentJobs(1, gocron.RescheduleMode)
		s1.Every(1).Hour().Do(service.NewChatEmail().HandleJobExpireToken)
		s1.StartAsync()
		defer s1.Clear()
	}
}
