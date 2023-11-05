package cmd

import (
	"context"
	"github.com/dyammarcano/application-manager/internal/service"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This service is a tool to generate the needed files
to quickly create a Cobra service.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.RegisterService("version", versionCall)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func versionCall(_ context.Context) error {
	service.AppVersion().String()
	return nil
}
