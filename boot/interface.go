package boot

import "github.com/hoangminhphuc/goph-chat/common/logger"


// Plug service in and play, unplug is easy
type Plugin func(*serviceHub)

type Service interface {
	Name() string
	InitFlags()
	Run() error
	Stop() <-chan error
}

// Services that need to be initialized before application runs
type InitService interface {
	Service
	Get() interface{}
}

// Services that run with the application runtime
type ActiveService interface {
	Service
}


// A service hub that handles all services
type ServiceHub interface {
	GetName() string
	GetLogger() logger.ZapLogger
	initFlags()
	parseFlags()
	Init() error
	Start() error
	Run() <-chan error
	Stop() error
	GetService(name string) (interface{}, bool)
	MustGetService(name string) interface{}
}
