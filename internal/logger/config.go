package logger

import (
	"fmt"
	"path/filepath"
)

const (
	defaultName        = "default"
	defaultInstance    = "local"
	defaultMaxFileSize = 100
	defaultMaxBackups  = 7
	defaultMaxAge      = 28
	defaultLocalTime   = true
	defaultCompress    = true
	defaultStdout      = true
)

type (
	Config struct {
		LogDir      string
		ServiceName string
		MaxFileSize int
		MaxAge      int
		MaxBackups  int
		LocalTime   bool
		Compress    bool
		Filename    string
		Stdout      bool
		Instance    string
		validated   bool
	}
)

func NewDefaultConfig() *Config {
	return &Config{
		LogDir:      "",
		ServiceName: defaultName,
		Instance:    defaultInstance,
		MaxFileSize: defaultMaxFileSize,
		MaxAge:      defaultMaxAge,
		MaxBackups:  defaultMaxBackups,
		LocalTime:   defaultLocalTime,
		Compress:    defaultCompress,
		Stdout:      defaultStdout,
		Filename:    fmt.Sprintf("%s-%s.log", defaultName, defaultInstance),
	}
}

func (c *Config) SetPath(logDir, name, instance string) error {
	if logDir == "" {
		c.LogDir = logDir
	}

	if name != "" {
		c.ServiceName = name
	}

	if instance != "" {
		c.Instance = instance
	}
	c.Stdout = false

	return c.Validate()
}

func (c *Config) SetRotation(maxFileSize int, maxAge int, maxBackups int, localTime, compress bool) {
	c.MaxAge = maxAge
	c.MaxFileSize = maxFileSize
	c.MaxBackups = maxBackups
	c.LocalTime = localTime
	c.Compress = compress
}

func (c *Config) Validate() error {
	c.validated = true

	if c.LogDir == "" {
		currPath, err := filepath.Abs(".")
		if err != nil {
			return err
		}

		c.LogDir = currPath
		c.Stdout = false
	}
	return nil
}
