package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "template-go",
	Short: "A brief description of your service",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your service. For example:

Cobra is a CLI library for Go that empowers applications.
This service is a tool to generate the needed files
to quickly create a Cobra service.`,
}

func init() {
	//rootCmd.Flags().StringVar(&service.CfgFile, "config", "", "config file")
	//rootCmd.Flags().StringVar(&service.CfgString, "config-string", "", "config string")
}
