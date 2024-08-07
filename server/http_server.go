package server

import (
	"net/http"
	"time"

	"github.com/tel4vn/fins-microservices/common/log"
	responsetime "github.com/tel4vn/fins-microservices/middleware/response"
	"github.com/tel4vn/fins-microservices/middleware/trace"
	"go.elastic.co/apm/module/apmgin/v2"

	"github.com/gin-gonic/gin"
)

const (
	serviceName = "test-service"
	version     = "v1.0.5"
)

func NewHTTPServer() *gin.Engine {
	// For Production
	// gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.MaxMultipartMemory = 100 << 20
	// APM
	engine.Use(apmgin.Middleware(engine))
	engine.Use(CORSMiddleware())
	engine.Use(allowOptionsMethod())
	engine.Use(trace.Handler)
	engine.Use(responsetime.Handler)
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": serviceName,
			"version": version,
			"time":    time.Now().Unix(),
		})
	})
	return engine
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Tenant-Id, X-Tenant-Uuid, X-Tenant-Name, X-User-Id, X-User-Level, X-User-Name")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")
		c.Next()
	}
}

func allowOptionsMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}

func Start(server *gin.Engine, port string) {
	v := make(chan struct{})
	go func() {
		if err := server.Run(":" + port); err != nil {
			log.Error("failed to start service", err)
			close(v)
		}
	}()
	log.Infof("service %v listening on port %v", serviceName, port)
	<-v
}
