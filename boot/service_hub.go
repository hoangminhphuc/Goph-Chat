package boot

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/facebookgo/flagenv"
	"github.com/hoangminhphuc/goph-chat/common/logger"
	"github.com/joho/godotenv"
)

type ServiceHub struct {
	name     			string
	services 		 	map[string]InitService
	Plugin  			[]Plugin
	signalChan   	chan os.Signal
	logger 	 			logger.ZapLogger
}

// No need to modify in the future so not returning pointer
func NewServiceHub(name string, Plugin ...Plugin) ServiceHub {
	service := ServiceHub{
		name:     name,
		services: map[string]InitService{},
		Plugin:   Plugin,
		logger: 	logger.NewZapLogger(),
	}



	// register services
	for _, p := range Plugin {
		p(&service)
	}

	service.initFlags()
	service.parseFlags()

	return service
}

func (s *ServiceHub) GetLogger() logger.ZapLogger {
	return s.logger
}

func (s *ServiceHub) initFlags() {
	for _, service := range s.services {
		service.InitFlags()
	}
}

func (s *ServiceHub) parseFlags() {
	err := godotenv.Load(".env")
	if err != nil {
		s.logger.Log.Error("Error loading env file ",err)
	}
	flagenv.Parse()
	flag.Parse()
	s.logger.Log.Info("Loading env file successfully")

}

// This function needs an instance of serviceHub, but when initialize service hub, 
// there are no instances so this method cannot be a method of serviceHub.
func RegisterPlugin(is InitService) Plugin {
	return func (s *ServiceHub) {
		if _, ok := s.services[is.Name()]; ok {
			log.Fatal("Service " + is.Name() + " already registered")
		}
		
		s.services[is.Name()] = is
	}
}

func (s *ServiceHub) Init() error {
	for _, sv := range s.services {
		if err := sv.Run(); err != nil {
			s.logger.Log.Error("Cannot initialize service ", sv.Name(), ". ", err.Error())
			return err
		}
	}
	return nil
}

func (s *ServiceHub) Start() error {
	signal.Notify(s.signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	c := s.Run()

	for {
		select {
		case err := <-c:
			if err != nil {
				s.logger.Log.Error(err.Error())
				s.Stop()
				return err
			}

		case sig := <-s.signalChan:
			s.logger.Log.Info(sig)
			switch sig {
				case syscall.SIGHUP:
					return nil
				default:
					s.Stop()
					return nil
			}
		}
	}
}

func (s *ServiceHub) Run() <-chan error {
	c := make(chan error, 1)

	for _, sv := range s.services {
		go func (s InitService) {
			c <- s.Run()
		}(sv)
	}
	return c
}

func (s *ServiceHub) Stop() error {
	s.logger.Log.Info("Stopping services...")
	stopChannel := make(chan error, len(s.services))

	for _, sv := range s.services {
		go func (s InitService) {
			stopChannel <- <-s.Stop()
		}(sv)
	}

	var errs []error

	// Wait for all services to stop
	for i := 0; i < len(s.services); i++ {
		if err := <- stopChannel; err != nil {
			s.logger.Log.Error("Failed to stop service: ", err)
			return err
		}
	}

	if len(errs) == 0 {
		s.logger.Log.Info("Service stopped successfully")
	}

	return nil
}