package cmd

import (
	"github.com/dyammarcano/template-go/internal/application"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "template-go",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute(manager *application.Application) {
	manager.Start(rootCmd)
}
