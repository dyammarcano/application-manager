package cmd

import (
	"fmt"
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
This service is a tool to generate the needed files
to quickly create a Cobra service.`,
	Run: func(cmd *cobra.Command, args []string) {
		service.RegisterService("client api", callClient)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

func callClient() error {
	fmt.Printf("client api: %v\n", time.Now().UTC())
	return nil
}
