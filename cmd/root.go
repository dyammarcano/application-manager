package cmd

import (
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
					logger.InfoAndPrint(fmt.Sprintf("simulate work: %s", ulid.Make()))
				}
			}
		})
	},
}

func Execute(version, commitHash, date string) {
	service.Execute(version, commitHash, date, rootCmd)
}

func init() {
	rootCmd.PersistentFlags().StringVar(&service.CfgFile, "config", "", "config file")
	rootCmd.PersistentFlags().StringVar(&service.CfgString, "config-str", "", "config string")
	rootCmd.PersistentFlags().StringVar(&service.LogsDir, "log-dir", "", "logs directory")
}
