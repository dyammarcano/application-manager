package application

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dyammarcano/application-manager/internal/algorithm/encoding"
	"github.com/dyammarcano/application-manager/internal/metadata"
	"github.com/dyammarcano/application-manager/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

var (
	AppVersion *metadata.Metadata
	CfgFile    string
	CfgString  string
)

type (
	Application struct {
		errChan  chan error
		ctx      context.Context
		ccf      context.CancelCauseFunc
		wGroup   sync.WaitGroup
		metadata *metadata.Metadata
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

		a.errChan <- fn(a.ctx)
	}()
}

// runServices executes all the services registered in the application
func (a *Application) runServices() {
	services := service.GetServices()
	if len(services) > 0 {
		a.chooseConfig()
		for name := range services {
			if runner, exist := services[name]; exist {
				a.executeInGoRoutine(runner)
			}
		}
		a.wGroup.Wait()
	}
}

func (a *Application) chooseConfig() {
	if CfgString != "" {
		if err := a.stringConfig(CfgString); err != nil {
			a.ccf(err)
		}
		return
	}

	if CfgFile == "" {
		a.loadConfigFileEnv()
		return
	}

	a.loadConfigFile()
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

func (a *Application) loadConfigFileEnv() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}
}

func (a *Application) loadConfigFile() {
	viper.SetConfigFile(CfgFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}
}

func (a *Application) stringConfig(data string) error {
	deserialized, err := encoding.Deserialize(data)
	if err != nil {
		return err
	}

	viper.SetConfigType("json")
	if err = viper.ReadConfig(bytes.NewBuffer([]byte(deserialized))); err != nil {
		return err
	}

	return nil
}
