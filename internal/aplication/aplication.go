package aplication

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	Client serviceName = iota
	Server
	Monitor
	Manager
)

var app *Application

func init() {
	app = newApp()
}

type (
	serviceName int

	Application struct {
		services map[serviceName]bool
		instance map[string]any
		errChan  chan error
		sigChan  chan os.Signal
		ctx      context.Context
		ccf      context.CancelCauseFunc
		wg       sync.WaitGroup
		mutex    sync.RWMutex
		cmd      *cobra.Command
	}
)

func newApp() *Application {
	app = &Application{
		services: make(map[serviceName]bool),
		errChan:  make(chan error),
		sigChan:  make(chan os.Signal),
	}

	signal.Notify(app.sigChan, syscall.SIGINT, syscall.SIGTERM)
	app.ctx, app.ccf = context.WithCancelCause(context.Background())

	return app
}

func (a *Application) addService(s serviceName) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.services[s] = true
}

func (a *Application) initConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		a.cmd.Println("Environment variables not found in: ", viper.ConfigFileUsed())
	}
}

func (a *Application) run(cmd *cobra.Command) {
	a.cmd = cmd
	a.errChan <- cmd.Execute()
}

func AddService(s serviceName) {
	app.addService(s)
}

func Run(cmd *cobra.Command) {
	app.run(cmd)
}
