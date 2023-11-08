package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	err := NewLoggerDefault()
	assert.Nil(t, err)

	Info("test")
	Error("test")
	Debug("test")
	Warn("test")
}

func TestNewLogger(t *testing.T) {
	cfg := NewDefaultConfig()
	err := cfg.SetPath("", "", "test")
	assert.Nil(t, err)

	err = NewLogger(cfg)
	assert.Nil(t, err)

	Info("test")
	Error("test")
	Debug("test")
}
