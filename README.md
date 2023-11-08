# cli template using custom cobra cli

## Todo


## Components

this example illustrates how to use the following components:

- service.AddFlag (to add flags to the service and the command)
- service.RegisterService (to register a service)
- service.Execute (to execute the service)


```go
package cmd

import (
    "context"
    "fmt"
    "github.com/dyammarcano/application-manager/internal/logger"
    "github.com/dyammarcano/application-manager/internal/service"
    "github.com/spf13/cobra"
    "time"
)

var rootCmd = &cobra.Command{
    ... rest of the basic code
    Run: func(cmd *cobra.Command, args []string) {
        service.RegisterService("main service", simulateWork)
    },
}

func Execute(version, commitHash, date string) {
    ctx := context.Background()
    service.Execute(ctx, version, commitHash, date, rootCmd)
}

func init() {
    // to add flags to be used by the service and the command
    service.AddFlag(rootCmd, "log-dir", "", "log file")
    service.AddFlag(rootCmd, "config", "", "config file")
    service.AddFlag(rootCmd, "config-string", "", "config string")
    service.AddFlag(rootCmd, "script", false, "script")
}

func simulateWork() error {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-service.Context().Done():
            return service.Context().Err()
        case <-ticker.C:
            logger.InfoAndPrint(fmt.Sprintf("simulate work: uuid: %s, ulid: %s, random: %s, time: %s",
                service.GetRandomValue("ulid"),
                service.GetRandomValue("uuid"),
                service.GetRandomValue("random"),
                time.Now().Format(time.RFC3339Nano),
            ))
        }
    }
}
