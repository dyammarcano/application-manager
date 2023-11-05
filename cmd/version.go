package cmd

import (
	"context"
	"github.com/dyammarcano/application-manager/internal/application"
	"github.com/dyammarcano/application-manager/internal/service"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.RegisterService("version", versionCall)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionCall(ctx context.Context) error {
	application.AppVersion.String()
	return nil
}
