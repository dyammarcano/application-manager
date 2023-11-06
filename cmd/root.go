package cmd

import (
	"context"
	"github.com/dyammarcano/application-manager/internal/service"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
	"log"
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
		service.RegisterService("all services", func(ctx context.Context) error {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-ticker.C:
					log.Print("simulate work: ", ulid.Make().String())
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
	rootCmd.PersistentFlags().StringVar(&service.CfgString, "config-string", "", "config string")
	rootCmd.PersistentFlags().StringVar(&service.LogsDir, "logs", "", "logs directory")
}
