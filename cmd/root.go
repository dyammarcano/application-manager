package cmd

import (
	"context"
	"fmt"
	"github.com/dyammarcano/application-manager/internal/command"
	"github.com/dyammarcano/application-manager/internal/logger"
	"github.com/dyammarcano/application-manager/internal/service"
	"github.com/spf13/cobra"
	"time"
)

var rootCmd = command.NewCommandBuilder("main").
	AddCommandShortMessage("A brief description of your service").
	AddCommandLongMessage(`A longer description that spans multiple lines and likely contains
examples and usage of using your service. For example:

Cobra is a CLI library for Go that empowers applications.
This service is a tool to generate the needed files
to quickly create a Cobra service.`).
	AddCommandRun(func(cmd *cobra.Command, args []string) {
		service.RegisterService("main service", simulateWork)
	}).
	AddCommandFlag("log-dir", "", "log file").
	AddCommandFlag("config", "", "config file").
	AddCommandFlag("config-string", "", "config string").
	AddCommandFlag("script", false, "script").
	Build()

func Execute(version, commitHash, date string) {
	ctx := context.Background()
	service.Execute(ctx, version, commitHash, date, rootCmd)
}

func simulateWork() error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-service.Context().Done():
			return service.Context().Err()
		case <-ticker.C:
			logger.Info(fmt.Sprintf("simulate work: uuid: %s, ulid: %s, time: %s",
				service.GetRandomValue("ulid"),
				service.GetRandomValue("uuid"),
				time.Now().Format(time.RFC3339Nano),
			))
		}
	}
}
