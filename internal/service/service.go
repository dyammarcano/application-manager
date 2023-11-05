package service

var s *Services

type (
	Runner func() error

	Services struct {
		services map[string]Runner
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
	s.services[serviceName] = runner
}

func GetServices() map[string]Runner {
	return s.services
}
