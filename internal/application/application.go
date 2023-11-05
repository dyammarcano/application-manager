package application

import (
	"context"
	"fmt"
	"github.com/dyammarcano/template-go/internal/metadata"
	"github.com/dyammarcano/template-go/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

var AppVersion *metadata.Metadata

type (
	Application struct {
		errChan       chan error
		ctx           context.Context
		ccf           context.CancelCauseFunc
		wGroup        sync.WaitGroup
		errorsRunning bool
		metadata      *metadata.Metadata
	}
)

// NewApplicationManager creates a new application manager
func NewApplicationManager(version, commitHash, date string) *Application {
	ctx, ccf := context.WithCancelCause(context.Background())

	app := &Application{
		errChan: make(chan error),
		wGroup:  sync.WaitGroup{},
		ctx:     ctx,
		ccf:     ccf,
		metadata: &metadata.Metadata{
			GoVersion:          runtime.Version(),
			ReleaseDate:        date,
			CommitHash:         commitHash,
			ApplicationVersion: version,
			Runtime: &metadata.Runtime{
				Arch: runtime.GOARCH,
				Goos: runtime.GOOS,
			},
		},
	}

	app.initConfig()
	AppVersion = app.metadata

	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		<-sigChan
		os.Exit(1)
	}()

	app.errorsHandler()

	return app
}

// initConfig reads in config file and ENV variables if set.
func (a *Application) initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Environment variables not found in: ", viper.ConfigFileUsed())
	}
}

// Start Execute uses the args (os.Args[1:] by default) and run services
func (a *Application) Start(cmd *cobra.Command) {
	a.errChan <- cmd.ExecuteContext(a.ctx)

	a.runServices()
}

// executeInGoRoutine executes a service in a go routine and returns the error in the error channel
func (a *Application) executeInGoRoutine(fn service.Runner) {
	a.wGroup.Add(1)

	go func() {
		defer a.wGroup.Done()

		a.errChan <- fn()
	}()
}

// runServices executes all the services registered in the application
func (a *Application) runServices() {
	services := service.GetServices()
	if len(services) > 0 {
		for name, runner := range services {
			fmt.Printf("Starting service: %s\n", name)
			a.executeInGoRoutine(runner)
		}
		a.wGroup.Wait()
	}
}

// errorsHandler handles the errors in the error channel
func (a *Application) errorsHandler() {
	go func() {
		for {
			select {
			case err := <-a.errChan:
				if err != nil {
					fmt.Println(err)
					a.ccf(err)
				}
			case <-a.ctx.Done():
				return
			}
		}
	}()
}
