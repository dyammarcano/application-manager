package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/caarlos0/log"
	"github.com/charmbracelet/lipgloss"
	"github.com/dyammarcano/application-manager/internal/algorithm/encoding"
	"github.com/dyammarcano/application-manager/internal/cache"
	"github.com/dyammarcano/application-manager/internal/logger"
	"github.com/dyammarcano/application-manager/internal/metadata"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/muesli/termenv"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/automaxprocs/maxprocs"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var ms *ManagerService

func init() {
	// enable colored output on github actions et al
	if os.Getenv("CI") != "" {
		lipgloss.SetColorProfile(termenv.TrueColor)
	}

	// automatically set GOMAXPROCS to match available CPUs.
	// GOMAXPROCS will be used as the default value for the --parallelism flag.
	if _, err := maxprocs.Set(); err != nil {
		log.WithError(err).Warn("failed to set GOMAXPROCS")
	}

	ms = &ManagerService{
		errChan:  make(chan error),
		wGroup:   sync.WaitGroup{},
		services: make(map[string]Runner),
		mutex:    sync.RWMutex{},
		v:        viper.New(),
	}
}

type (
	Runner func() error

	ManagerService struct {
		errChan   chan error
		ctx       context.Context
		causeFunc context.CancelCauseFunc
		wGroup    sync.WaitGroup
		metadata  *metadata.Metadata
		services  map[string]Runner
		mutex     sync.RWMutex
		v3c       *cache.V3Cache
		v         *viper.Viper
	}
)

// setup initializes the service manager
func setup(ctx context.Context, version, commitHash, date string) {
	ms.ctx, ms.causeFunc = context.WithCancelCause(ctx)
	setupOsExitHandler(ms.ctx)
	ms.metadata = initMetadata(version, commitHash, date)
	ms.errorsHandler()
}

// AddFlag adds a flag to the service manager, it also binds the flag to the viper instance
func AddFlag(cmd *cobra.Command, name string, defaultValue any, description string) {
	switch v := defaultValue.(type) {
	case bool:
		cmd.PersistentFlags().Bool(name, v, description)
	case string:
		cmd.PersistentFlags().String(name, v, description)
	case int, int8, int16, int32, int64:
		cmd.PersistentFlags().Int64(name, v.(int64), description)
	default:
		fmt.Printf("Invalid type: %s\n", v)
		os.Exit(1)
	}

	if err := ms.v.BindPFlag(name, cmd.PersistentFlags().Lookup(name)); err != nil {
		cmd.Printf("Error binding flag: %s\n", err)
		os.Exit(1)
	}
}

// GetValue returns the flag value
func GetValue(name string) any {
	return ms.v.Get(name)
}

// SetValue sets the flag value
func SetValue(name string, value any) {
	ms.v.Set(name, value)
}

// setupOsExitHandler sets up the os exit handler
func setupOsExitHandler(ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-sigChan:
			log.Info("receiving signal to gracefully exiting")
			os.Exit(1)
		case <-ctx.Done():
			return
		}
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
	if ms == nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Context returns the global context
func Context() context.Context {
	return ms.ctx
}

// Execute creates a new service manager
func Execute(ctx context.Context, version, commitHash, date string, cmd *cobra.Command) {
	setup(ctx, version, commitHash, date)
	ms.errChan <- cmd.ExecuteContext(ms.ctx)

	if ms.v.GetBool("script") == true {
		ms.generateScript()
	}

	ms.runServices()
}

// AppVersion returns the service version
func AppVersion() *metadata.Metadata {
	errAndExit("service instance is not initialized")
	return ms.metadata
}

// RegisterService adds a service to the service to be executed
func RegisterService(serviceName string, runner Runner) {
	errAndExit("service instance is not initialized")
	ms.registerService(serviceName, runner)
}

// GetRandomValue returns a random guid like ulid, uuid, string, etc
func GetRandomValue(name string) string {
	switch name {
	case "ulid":
		return ulid.Make().String()
	case "uuid":
		return uuid.New().String()
	case "random":
		h := sha256.New()
		h.Write([]byte(ulid.Make().String()))
		return fmt.Sprintf("%x", h.Sum(nil))
	default:
		return ""
	}
}

// RegisterService adds a service to the service to be executed
func (a *ManagerService) registerService(serviceName string, runner Runner) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.services[serviceName] = runner
	log.Infof("[stage 0] service registered: [%s]", serviceName)
}

// executeInGoRoutine executes a service in a go routine and returns the error in the error channel
func (a *ManagerService) executeInGoRoutine(fn Runner) {
	a.wGroup.Add(1)

	go func() {
		defer a.wGroup.Done()

		log.Infof("ready=1")
		a.errChan <- fn()
	}()
}

// runServices executes all the services registered in the service
func (a *ManagerService) runServices() {
	if len(a.services) > 0 {
		a.chooseConfig()
		a.setupLogger()
		for name := range a.services {
			if runner, exist := a.services[name]; exist {
				log.Infof("[stage 1] starting service: [%s]", name)
				a.executeInGoRoutine(runner)
			}
		}
		a.wGroup.Wait()
	}
}

// chooseConfig chooses the config to be used
func (a *ManagerService) chooseConfig() {
	go a.watchConfig()

	configStr := ms.v.GetString("config-string")

	if configStr != "" {
		a.stringConfig(configStr)
		return
	}

	cfgFile := ms.v.GetString("config")

	if cfgFile == "" {
		a.loadConfigFileEnv()
		return
	}

	a.loadConfigFile(cfgFile)
}

// setupLogger check if logger are set in config or by commanf flag
func (a *ManagerService) setupLogger() {
	logPath := a.v.GetString("log-dir")

	if logPath == "" {
		if err := logger.NewLoggerDefault(); err != nil {
			a.causeFunc(err)
		}
		return
	}

	cfg := logger.NewDefaultConfig()
	if err := cfg.SetPath(logPath, "", ""); err != nil {
		a.causeFunc(err)
	}

	if err := logger.NewLogger(cfg); err != nil {
		a.causeFunc(err)
	}
}

// errorsHandler handles the errors in the error channel
func (a *ManagerService) errorsHandler() {
	go func() {
		for {
			select {
			case err := <-a.errChan:
				if err != nil {
					fmt.Println(err)
					a.causeFunc(err)
				}
			case <-a.ctx.Done():
				return
			}
		}
	}()
}

// loadConfigFileEnv loads the config file from the environment
func (a *ManagerService) loadConfigFileEnv() {
	ms.v.AddConfigPath(".")
	ms.v.SetConfigName("app")
	ms.v.SetConfigType("env")
	ms.v.AutomaticEnv()

	if err := ms.v.ReadInConfig(); err == nil {
		log.Infof("[stage 0] using config file: %s", ms.v.ConfigFileUsed())
	}
}

// loadConfigFile loads the config file from the file system
func (a *ManagerService) loadConfigFile(cfgFile string) {
	ms.v.SetConfigFile(cfgFile)
	ms.v.AutomaticEnv()

	if err := ms.v.ReadInConfig(); err == nil {
		log.Infof("[stage 0] using config file: %s", ms.v.ConfigFileUsed())
	}
}

// stringConfig loads the config from a string
func (a *ManagerService) stringConfig(data string) {
	deserialized, err := encoding.Deserialize(data)
	if err != nil {
		a.causeFunc(err)
	}

	ms.v.SetConfigType("json")
	if err = ms.v.ReadConfig(bytes.NewBuffer([]byte(deserialized))); err != nil {
		a.causeFunc(err)
	}
}

// watchConfig watches the config file for changes
func (a *ManagerService) watchConfig() {
	<-time.After(5 * time.Second)

	if ms.v.GetString("config-string") != "" {
		ms.v.WatchConfig()
		ms.v.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("[stage 0] config file changed:", e.Name)
		})
	}
}

// generateScript generates the script for the service
func (a *ManagerService) generateScript() {
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

func (a *ManagerService) checkForUpdates() {
	// check for updates
	// download updates
	// install updates
}

func (a *ManagerService) validateUpdate() {
	// validate updates
	// validate config
}

func (a *ManagerService) downloadUpdate() {
	// download updates
	// download config
}
