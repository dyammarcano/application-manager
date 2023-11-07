package cmd

import (
	"context"
	"fmt"
	"github.com/dyammarcano/application-manager/internal/logger"
	"github.com/dyammarcano/application-manager/internal/service"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "template-go",
	Short: "A brief description of your service",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your service. For example:

Cobra is a CLI library for Go that empowers applications.
This service is a tool to generate the needed files
to quickly create a Cobra service.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.RegisterService("main service", func() error {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-service.Context().Done():
					return service.Context().Err()
				case <-ticker.C:
					logger.InfoAndPrint(fmt.Sprintf("simulate work: %s, %v", ulid.Make(), service.GetValue("picture")))
				}
			}
		})
	},
}

func Execute(version, commitHash, date string) {
	ctx := context.Background()
	service.Execute(ctx, version, commitHash, date, rootCmd)
}

func init() {
	service.AddFlag(rootCmd, "log-dir", "", "log file")
	service.AddFlag(rootCmd, "config", "", "config file")
	service.AddFlag(rootCmd, "config-string", "", "config string")
	service.AddFlag(rootCmd, "script", false, "script")
}
