package service

import (
	"context"
	"sync"
)

var s *Services

type (
	Runner func(ctx context.Context) error

	Services struct {
		services map[string]Runner
		mutex    sync.Mutex
	}
)

func init() {
	s = &Services{
		services: make(map[string]Runner),
	}
}

// RegisterService adds a service to the application to be executed
func RegisterService(serviceName string, runner Runner) {
	s.registerService(serviceName, runner)
}

func (s *Services) registerService(serviceName string, runner Runner) {
	defer s.mutex.Unlock()
	s.mutex.Lock()

	s.services[serviceName] = runner
}

func GetServices() map[string]Runner {
	defer s.mutex.Unlock()
	s.mutex.Lock()

	return s.services
}
