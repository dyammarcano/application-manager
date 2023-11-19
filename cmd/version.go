package cmd

import (
	"github.com/dyammarcano/application-manager/internal/command"
	"github.com/dyammarcano/application-manager/internal/service"
	"github.com/spf13/cobra"
)

var versionCmd = command.NewCommandBuilder("version").
	AddCommandShortMessage("A brief description of your command").
	AddCommandLongMessage(`A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This service is a tool to generate the needed files
to quickly create a Cobra service.`).
	AddCommandRun(func(cmd *cobra.Command, args []string) {
		service.RegisterService("version", versionCall)
	}).
	Build()

func init() {
	rootCmd.AddCommand(versionCmd)
}

func versionCall() error {
	service.AppVersion().String()
	return nil
}
