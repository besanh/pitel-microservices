package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/tel4vn/fins-microservices/common/env"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/internal/redis"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
	logLevel := env.GetStringENV("LOG_LEVEL", "error")
	logFile := "tmp/console.log"
	log.InitLogger(logLevel, logFile)
}

func initRedis() {
	var err error
	log.Info(env.GetIntENV("REDIS_DATABASE", 0))
	if redis.Redis, err = redis.NewRedis(redis.Config{
		Addr:         env.GetStringENV("REDIS_ADDRESS", "localhost:6379"),
		Password:     env.GetStringENV("REDIS_PASSWORD", ""),
		DB:           env.GetIntENV("REDIS_DATABASE", 0),
		PoolSize:     30,
		PoolTimeout:  20,
		IdleTimeout:  10,
		ReadTimeout:  20,
		WriteTimeout: 15,
	}); err != nil {
		panic(err)
	}
}

func Execute() {
	var rootCmd = &cobra.Command{Use: "bss-inbox-marketing"}
	rootCmd.AddCommand(InitCMDAPI())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
