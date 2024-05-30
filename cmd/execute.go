package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tel4vn/fins-microservices/common/log"
)

func Execute() {
	var rootCmd = cobra.Command{Use: "chat-service"}
	rootCmd.AddCommand(cmdMain)
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
