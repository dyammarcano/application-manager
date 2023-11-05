package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/dyammarcano/template-go/internal/helpers"
	"github.com/dyammarcano/template-go/internal/service"
	"io"
	"net/http"
	"runtime/trace"

	"github.com/spf13/cobra"
)

type (
	ExternalAPIResponse struct {
		Args    map[string]interface{} `json:"args"`
		Headers map[string]interface{} `json:"headers"`
		Origin  string                 `json:"origin"`
		URL     string                 `json:"url"`
	}

	ExternalAPIResponseError struct {
		Error string `json:"error"`
	}
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
			apiUrl, err := helpers.GetString("url")
			if err != nil {
				return err
			}

			request, err := http.NewRequestWithContext(cmd.Context(), "GET", apiUrl, nil)
			if err != nil {
				return err
			}

			client := &http.Client{}

			defer trace.StartRegion(cmd.Context(), "client.Do").End()
			response, err := client.Do(request)
			if err != nil {
				return err
			}

			defer trace.StartRegion(cmd.Context(), "response.Body.Close").End()
			data, err := io.ReadAll(response.Body)
			if err != nil {
				return err
			}

			if response.StatusCode != 200 {
				var errResponse ExternalAPIResponseError
				if err = json.Unmarshal(data, &errResponse); err != nil {
					return err
				}
				fmt.Fprint(cmd.OutOrStdout(), errResponse.Error)
				return nil
			}

			var apiResponse ExternalAPIResponse

			defer trace.StartRegion(cmd.Context(), "json.Unmarshal").End()
			if err = json.Unmarshal(data, &apiResponse); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "your public ip: %s", apiResponse.Origin)
			return nil
		})
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
