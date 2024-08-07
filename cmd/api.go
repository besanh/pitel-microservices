package cmd

import (
	"fmt"
	"net/http"

	"github.com/cardinalby/hureg"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	v1 "github.com/tel4vn/fins-microservices/api/v1"
	"github.com/tel4vn/fins-microservices/common/cache"
	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/response"
	"github.com/tel4vn/fins-microservices/internal/goauth"
	"github.com/tel4vn/fins-microservices/internal/redis"
	"github.com/tel4vn/fins-microservices/internal/sqlclient"
	authMdw "github.com/tel4vn/fins-microservices/middleware/auth"
	"github.com/tel4vn/fins-microservices/server"
	"github.com/tel4vn/fins-microservices/service"
)

func InitCMDAPI() *cobra.Command {
	var cmdAPI = &cobra.Command{
		Use: "api",
		Run: func(cmd *cobra.Command, args []string) {
			RunAPI(cmd)
		},
	}
	return cmdAPI
}

func RunAPI(cmd *cobra.Command) {
	log.Debug("starting api service")
	// _, err := cmd.Flags().GetString("service")
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// switch sv {
	// case "bss-inbox-marketing-service":
	RunAPIBssInboxMarketing()
	// default:
	// 	log.Error("unknown service")
	// }
}

func RunAPIBssInboxMarketing() {
	// Init redis
	log.Debug("Initializing redis")
	initRedis()

	// Init mem cache
	log.Debug("Initializing mem cache")
	cache.MCache = cache.NewMemCache()

	// Init redis cache
	log.Debug("Initializing redis cache")
	cache.RCache = cache.NewRedisCache(redis.Redis.GetClient())

	// Init middlewares
	log.Debug("Initializing middleware")
	goauth.GoAuthClient = goauth.NewGoAuth(redis.Redis.GetClient())
	authMdw.SECRET_TOKEN = env.GetStringENV("SECRET_TOKEN", "ZXllRmluUy5TZWNyZXQyMDIzLmIyazBXRGxwY0ROcVZXazVNazFaVW1aS2NsSldNbA==")

	// DB Pool
	dbConfig := sqlclient.SqlConfig{
		Host:         env.GetStringENV("DB_HOST", "localhost"),
		Database:     env.GetStringENV("DB_DATABASE", "bss-inbox-marketing"),
		Username:     env.GetStringENV("DB_USERNAME", "admin"),
		Password:     env.GetStringENV("DB_PASSWORD", "password"),
		Port:         env.GetIntENV("DB_PORT", 5432),
		DialTimeout:  20,
		ReadTimeout:  30,
		WriteTimeout: 30,
		Timeout:      30,
		PoolSize:     10,
		MaxOpenConns: 20,
		MaxIdleConns: 10,
		Driver:       sqlclient.POSTGRESQL,
	}

	db, err := sqlclient.NewSqlClient(dbConfig)
	if err != nil {
		log.Error(err)
		return
	}

	service.InitDBConnection(db)

	router := server.NewHTTPServer()
	humaAPI := humagin.New(router, huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "1.0.0",
			Info: &huma.Info{
				Title:       "BSS Inbox Marketing APIs",
				Version:     "1.0.0",
				Description: "BSS Inbox Marketing APIs",
				Contact: &huma.Contact{
					Name:  "TEL4VN-Pitel",
					URL:   "https://www.pitel.vn/",
					Email: "innovation@tel4vn.com",
				},
			},
			Components: &huma.Components{
				SecuritySchemes: map[string]*huma.SecurityScheme{
					"bssAuth": {
						Type:         "http",
						Scheme:       "bearer",
						In:           "header",
						Description:  "Authorization header using the Bearer scheme. Example: \"Authorization: Bearer {token}\"",
						BearerFormat: "Token String",
						Name:         "Authorization",
					},
				},
			},
			Servers: []*huma.Server{
				{
					URL:         "https://api.dev.fins.vn",
					Description: "Development Environment",
					Variables:   map[string]*huma.ServerVariable{},
				},
				{
					URL:         "https://api.uat.dev.fins.vn",
					Description: "UAT Environment",
					Variables:   map[string]*huma.ServerVariable{},
				},
				{
					URL:         "https://api.dev.fins.vn",
					Description: "Production Environment",
					Variables:   map[string]*huma.ServerVariable{},
				},
				{
					URL:         "http://localhost:8009",
					Description: "Local Environment",
					Variables:   map[string]*huma.ServerVariable{},
				},
			},
		},
		OpenAPIPath:   "/docs/openapi",
		DocsPath:      "",
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
	})
	api := hureg.NewAPIGen(humaAPI)
	api.GetHumaAPI().UseMiddleware(authMdw.NewAuthMiddleware(api))

	v1.RegisterAPIAuth(api)

	log.Debug("Starting api server")
	response.NewHumaError()
	port := env.GetIntENV("PORT", 8000)

	router.GET("/docs/api-document", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
		<!doctype html>
		<html>
			<head>
				<title>BSS Inbox Marketing APIs</title>
				<meta charset="utf-8" />
				<meta
				name="viewport"
				content="width=device-width, initial-scale=1" />
			</head>
			<body>
				<script
				id="api-reference"
				data-url="/collection/openapi.json"></script>
				<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
			</body>
		</html>
		`))
	})

	// Init cron job

	server.Start(router, fmt.Sprintf("%d", port))
}
