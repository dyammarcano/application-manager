package cmd

import (
	"fmt"
	"github.com/dyammarcano/application-manager/internal/command"
	"github.com/dyammarcano/application-manager/internal/service"
	"github.com/spf13/cobra"
	"time"
)

var clientCmd = command.NewCommandBuilder("client").
	AddCommandShortMessage("A brief description of your command").
	AddCommandLongMessage(`A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This service is a tool to generate the needed files
to quickly create a Cobra service.`).
	AddCommandRun(func(cmd *cobra.Command, args []string) {
		service.RegisterService("client api", callClient)
	}).
	Build()

func init() {
	rootCmd.AddCommand(clientCmd)
}

func callClient() error {
	fmt.Printf("client api: %v\n", time.Now().UTC())
	return nil
}
