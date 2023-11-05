package cmd

import (
	"github.com/dyammarcano/application-manager/internal/service"
	"github.com/spf13/cobra"
	"time"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.RegisterService("client api", func() error {
			cmd.Printf("client api: %v\n", time.Now().UTC())
			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
