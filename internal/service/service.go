package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/dyammarcano/application-manager/internal/algorithm/encoding"
	"github.com/dyammarcano/application-manager/internal/cache"
	"github.com/dyammarcano/application-manager/internal/logger"
	"github.com/dyammarcano/application-manager/internal/metadata"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const (
	cfgEnv cfgSource = iota
	cfgFile
	stringCfg
)

var (
	svc        *Application
	CfgFile    string
	CfgString  string
	ScriptFlag bool
	LogsDir    string
)

type (
	cfgSource int

	Runner func() error

	Application struct {
		errChan   chan error
		ctx       context.Context
		ccf       context.CancelCauseFunc
		wGroup    sync.WaitGroup
		metadata  *metadata.Metadata
		services  map[string]Runner
		mutex     sync.RWMutex
		cfgSource cfgSource
		zlog      *logger.ZapLogger
		v3c       *cache.V3Cache
	}
)

// setup initializes the service manager
func setup(version, commitHash, date string) {
	ctx, ccf := context.WithCancelCause(context.Background())
	setupOsExitHandler()

	svc = &Application{
		errChan:  make(chan error),
		wGroup:   sync.WaitGroup{},
		ctx:      ctx,
		ccf:      ccf,
		services: make(map[string]Runner),
		mutex:    sync.RWMutex{},
		metadata: initMetadata(version, commitHash, date),
	}

	svc.errorsHandler()
}

// setupOsExitHandler sets up the os exit handler
func setupOsExitHandler() {
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		<-sigChan
		log.Printf("receiving signal to gracefully exiting")
		os.Exit(1)
	}()
}

// initMetadata initializes the metadata
func initMetadata(version, commitHash, date string) *metadata.Metadata {
	return &metadata.Metadata{
		GoVersion:          runtime.Version(),
		ReleaseDate:        date,
		CommitHash:         commitHash,
		ApplicationVersion: version,
		Runtime: &metadata.Runtime{
			Arch: runtime.GOARCH,
			Goos: runtime.GOOS,
			Pid:  os.Getpid(),
			PPid: os.Getppid(),
		},
	}
}

// errAndExit prints the error and exits the service
func errAndExit(err any) {
	if svc == nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Logger() *logger.ZapLogger {
	return svc.zlog
}

func Context() context.Context {
	return svc.ctx
}

// Execute creates a new service manager
func Execute(version, commitHash, date string, cmd *cobra.Command) {
	setup(version, commitHash, date)
	svc.errChan <- cmd.ExecuteContext(svc.ctx)

	if ScriptFlag == true {
		svc.generateScript()
	}

	svc.runServices()
}

// AppVersion returns the service version
func AppVersion() *metadata.Metadata {
	errAndExit("service instance is not initialized")
	return svc.metadata
}

// RegisterService adds a service to the service to be executed
func RegisterService(serviceName string, runner Runner) {
	errAndExit("service instance is not initialized")
	svc.registerService(serviceName, runner)
}

// GetContext returns global context
func GetContext() context.Context {
	return svc.ctx
}

// RegisterService adds a service to the service to be executed
func (a *Application) registerService(serviceName string, runner Runner) {
	defer a.mutex.Unlock()
	a.mutex.Lock()

	a.services[serviceName] = runner
	log.Printf("service registered: [%s]\n", serviceName)
}

// executeInGoRoutine executes a service in a go routine and returns the error in the error channel
func (a *Application) executeInGoRoutine(fn Runner) {
	a.wGroup.Add(1)

	go func() {
		defer a.wGroup.Done()

		a.errChan <- fn()
	}()
}

// runServices executes all the services registered in the service
func (a *Application) runServices() {
	if len(a.services) > 0 {
		a.chooseConfig()
		a.setupLogger()
		for name := range a.services {
			if runner, exist := a.services[name]; exist {
				log.Printf("starting service: [%s]\n", name)
				a.executeInGoRoutine(runner)
			}
		}
		a.wGroup.Wait()
	}
}

// chooseConfig chooses the config to be used
func (a *Application) chooseConfig() {
	go a.watchConfig()

	if CfgString != "" {
		a.cfgSource = stringCfg
		a.stringConfig(CfgString)
		return
	}

	if CfgFile == "" {
		a.cfgSource = cfgEnv
		a.loadConfigFileEnv()
		return
	}

	a.cfgSource = cfgFile
	a.loadConfigFile()
}

// setupLogger check if logger are set in config or by commanf flag
func (a *Application) setupLogger() {
	elogPath := viper.GetString("logPath")
	if LogsDir != "" || elogPath != "" {
		cfg := logger.NewDefaultConfig()
		if err := cfg.SetPath(LogsDir, "", ""); err != nil {
			a.ccf(err)
		}

		ll, err := logger.NewLogger(cfg)
		if err != nil {
			a.ccf(err)
		}
		a.zlog = ll
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

// loadConfigFileEnv loads the config file from the environment
func (a *Application) loadConfigFileEnv() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}
}

// loadConfigFile loads the config file from the file system
func (a *Application) loadConfigFile() {
	viper.SetConfigFile(CfgFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}
}

// stringConfig loads the config from a string
func (a *Application) stringConfig(data string) {
	deserialized, err := encoding.Deserialize(data)
	if err != nil {
		a.ccf(err)
	}

	viper.SetConfigType("json")
	if err = viper.ReadConfig(bytes.NewBuffer([]byte(deserialized))); err != nil {
		a.ccf(err)
	}
}

// watchConfig watches the config file for changes
func (a *Application) watchConfig() {
	<-time.After(5 * time.Second)

	if a.cfgSource != stringCfg {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("config file changed:", e.Name)
		})
	}
}

// generateScript generates the script for the service
func (a *Application) generateScript() {
	for name := range a.services {
		if _, exist := a.services[name]; exist {
			fmt.Println("generating script for service: ", name)
			// generate script
			// get service name
			// get service config and serialize it
			// output script to file
		}
	}
}

func (a *Application) checkForUpdates() {
	// check for updates
	// download updates
	// install updates
}

func (a *Application) validateUpdate() {
	// validate updates
	// validate config
}

func (a *Application) downloadUpdate() {
	// download updates
	// download config
}

//func (a *Application) loadScript() {
//	// load script
//	// deserialize config
//	// load config
//}
//
//func (a *Application) runScript() {
//	// run script
//	// run service
//}
//
//func (a *Application) stopScript() {
//	// stop script
//	// stop service
//}
