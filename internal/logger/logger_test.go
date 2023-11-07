package logger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	ll, err := NewLoggerDefault()
	assert.Nil(t, err)

	ll.Info("test")
	ll.Error("test")
}

func TestNewLogger(t *testing.T) {
	cfg := NewDefaultConfig()
	err := cfg.SetPath("", "", "test")
	assert.Nil(t, err)

	ll, err := NewLogger(cfg)
	assert.Nil(t, err)

	ll.Info("test")
	ll.Error("test")
}
