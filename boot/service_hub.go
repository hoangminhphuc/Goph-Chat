package boot

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/facebookgo/flagenv"
	"github.com/hoangminhphuc/goph-chat/common"
	"github.com/hoangminhphuc/goph-chat/common/logger"
	rt "github.com/hoangminhphuc/goph-chat/internal/router"
	"github.com/hoangminhphuc/goph-chat/internal/server/websocket"
	"github.com/hoangminhphuc/goph-chat/module/room/model"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type serviceHub struct {
	name     			string
	runtimeService []RuntimeService
	initServices 	map[string]InitService
	Plugin  			[]Plugin
	httpServer 		*rt.HTTPServer
	signalChan   	chan os.Signal
	logger 	 			logger.ZapLogger
}

// No need to modify in the future so not returning pointer
func NewServiceHub(name string, Plugin ...Plugin) ServiceHub {
	service := &serviceHub{
		name:     		name,
		initServices: 		map[string]InitService{},
		Plugin:   		Plugin,
		signalChan:   make(chan os.Signal, 1),
		logger: 			logger.NewZapLogger(),
	}



	// register services
	for _, p := range Plugin {
		p(service)
	}

	httpServer := rt.NewHTTPServer()
	service.httpServer = httpServer
	service.runtimeService = append(service.runtimeService, httpServer)

	ws := websocket.NewWebSocketServer()
	service.runtimeService = append(service.runtimeService, ws)

	service.initFlags()
	service.parseFlags()

	return service
}

func (s *serviceHub) GetName() string {
	return s.name
}
func (s *serviceHub) GetLogger() logger.ZapLogger {
	return s.logger
}

func (s *serviceHub) GetHTTPServer() *rt.HTTPServer {
	return s.httpServer
}

func (s *serviceHub) GetWSServer() *websocket.WebSocketServer {
	for _, as := range s.runtimeService {
		if as.Name() == "websocket" {
			return as.(*websocket.WebSocketServer)
		}
	}

	return nil
}

// ! Will be refactoring later on
func (s *serviceHub) InitializePools(ws *websocket.WebSocketServer) {
	db := s.MustGetService(common.PluginDBMain).(*gorm.DB)

	var rooms []model.Room
	if err := db.Find(&rooms).Error; err != nil {
			s.logger.Log.Error(err.Error())
	}

	for _, r := range rooms {
			pool := websocket.NewPool(r.ID)
			ws.Rooms[r.ID] = pool
	}
}

func (s *serviceHub) initFlags() {
	for _, is := range s.initServices {
		is.InitFlags()
	}

	for _, as := range s.runtimeService {
		as.InitFlags()
	}
}

func (s *serviceHub) parseFlags() {
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
	return func (s *serviceHub) {
		if _, ok := s.initServices[is.Name()]; ok {
			log.Fatal("Service " + is.Name() + " already registered")
		}
		
		s.initServices[is.Name()] = is
	}
}


// Initialize all services, no need for listening to error
func (s *serviceHub) Init() error {
	for _, sv := range s.initServices {
		if err := sv.Run(); err != nil {
			s.logger.Log.Error("Cannot initialize service ", sv.Name(), ". ", err.Error())
			return err
		}
	}
	return nil
}

func (s *serviceHub) Start() error {
	/* 
		Whenever the OS sends a signal (Ctrl+C or kill), 
		the programm won't handle it (which is by default). 
		Instead, send that signal into s.signalChan so that we can handle it
	*/
	signal.Notify(s.signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	c := s.Run()
	s.logger.Log.Info("Starting runtime services...")

	for {
		select {
		case err := <-c:
			if err != nil {
				s.logger.Log.Error(err.Error())
				s.Stop()
				return err
			}

		case sig := <-s.signalChan:
			s.logger.Log.Info("Received signal: ", sig)
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

func (s *serviceHub) Run() <-chan error {
	c := make(chan error, 1)

	for _, sv := range s.runtimeService {
		go func (s RuntimeService) {
			c <- s.Run()
		}(sv)
	}
	return c
}

func (s *serviceHub) Stop() error {
	s.logger.Log.Info("Stopping services...")

	var (
		lenService = len(s.runtimeService)+len(s.initServices)
		stopChannel = make(chan error, lenService)
	)


	// Even no error (nil) when stopping, still sends to stopChannel 

	for _, is := range s.initServices {
		go func(ins InitService) {
			stopChannel <- <-ins.Stop()
		}(is)
	}


	for _, as := range s.runtimeService {
		go func(acs RuntimeService) {
			stopChannel <- <-acs.Stop()
		}(as)
	}


	for i := 0; i < lenService; i++ {
		if err := <- stopChannel; err != nil {
			s.logger.Log.Error("Failed to stop service: ", err)
			return err
		}
	}
	
	s.logger.Log.Info("Service stopped successfully")
	return nil
}


func (s *serviceHub) GetService(name string) (interface{}, bool) {
	is, ok := s.initServices[name]

	if !ok {
		return nil, ok
	}

	return is.Get(), true
}

func (s *serviceHub) MustGetService(name string) interface{} {
	sv, ok := s.GetService(name)

	if !ok {
		s.logger.Log.Fatal("Service " + name + " not found")
	}

	return sv
}

func (s *serviceHub) GetEnvValue(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		s.logger.Log.Error("Environment variable not found: ", key)
		return ""
	}
	return value
}