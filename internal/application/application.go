package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

var app *Application

type (
	ServiceRunner func() error

	Application struct {
		services      map[string]ServiceRunner
		errChan       chan error
		sigChan       chan os.Signal
		ctx           context.Context
		ccf           context.CancelCauseFunc
		wGroup        sync.WaitGroup
		metadata      Metadata
		errorsRunning bool
	}

	Metadata struct {
		ApplicationVersion string  `json:"application_version"`
		CommitHash         string  `json:"commit_hash"`
		GoVersion          string  `json:"go_version"`
		ReleaseDate        string  `json:"release_date"`
		CommitTag          string  `json:"commit_tag"`
		Runtime            Runtime `json:"runtime"`
	}

	Runtime struct {
		Arch string `json:"arch"`
		Goos string `json:"goos"`
	}
)

func newApp(version, commitHash, date string) {
	if app == nil {
		ctx, ccf := context.WithCancelCause(context.Background())

		app = &Application{
			services: make(map[string]ServiceRunner),
			errChan:  make(chan error),
			sigChan:  make(chan os.Signal),
			ctx:      ctx,
			ccf:      ccf,
			metadata: Metadata{
				GoVersion:          runtime.Version(),
				ReleaseDate:        date,
				CommitHash:         commitHash,
				ApplicationVersion: version,
				Runtime: Runtime{
					Arch: runtime.GOARCH,
					Goos: runtime.GOOS,
				},
			},
		}
		signal.Notify(app.sigChan, syscall.SIGINT, syscall.SIGTERM)
	}
}

func (a *Application) initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Environment variables not found in: ", viper.ConfigFileUsed())
	}
}

func (a *Application) registerService(serviceName string, runner ServiceRunner) {
	a.services[serviceName] = runner
}

func (a *Application) execute(cmd *cobra.Command) {
	a.errChan <- cmd.ExecuteContext(a.ctx)

	a.initConfig()
	a.errChan <- a.runServices()
}

func (a *Application) waitGroupsToFinish() {
	a.wGroup.Wait()
}

func (a *Application) executeInGoRoutine(fn ServiceRunner) {
	a.wGroup.Add(1)

	go func() {
		defer a.wGroup.Done()

		a.errChan <- fn()
	}()
}

func (a *Application) runServices() error {
	if len(a.services) == 0 {
		return errors.New("no services to run")
	}

	for name, runner := range a.services {
		fmt.Printf("Starting service: %s\n", name)
		a.executeInGoRoutine(runner)
	}

	a.waitGroupsToFinish()
	return nil
}

func (a *Application) errors() {
	if a.errorsRunning {
		return
	}

	go func() {
		a.errorsRunning = true
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

func (a *Application) version() (string, error) {
	d, err := json.Marshal(a.metadata)
	if err != nil {
		return "", err
	}
	return string(d) + "\n", nil
}

func AppVersion() (string, error) {
	if app == nil {
		return "", fmt.Errorf("application not initialized")
	}

	return app.version()
}

func RegisterService(serviceName string, runner ServiceRunner) {
	app.registerService(serviceName, runner)
}

func Init(version, commitHash, date string) {
	newApp(version, commitHash, date)

	app.errors()
}

func ExecuteCommand(cmd *cobra.Command) {
	app.execute(cmd)
}
