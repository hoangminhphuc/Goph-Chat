package router

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hoangminhphuc/goph-chat/common/logger"
)

var (
	defaultHost         = "localhost"
	defaultPort         = "8080"
	defaultReadTimeout  = 15 * time.Second
	defaultWriteTimeout = 15 * time.Second
	defaultIdleTimeout  = 60 * time.Second
)

type Config struct {
	Host 						string
  Port 						string
	ReadTimeout 		time.Duration
	WriteTimeout 		time.Duration
	IdleTimeout     time.Duration
}

type HTTPServer struct {
	name 			string
	router 		*gin.Engine
	server    *http.Server
	logger 		logger.ZapLogger
	Config
}

func NewHTTPServer() *HTTPServer {
	srv := &HTTPServer {
		name: 		"gin",
		router: 	gin.Default(),
		logger: 	logger.NewZapLogger(),
	}
	return srv
}

func (rs *HTTPServer) GetName() string {
	return rs.name
}

func (rs *HTTPServer) GetRouter() *gin.Engine {
	return rs.router
}

func (rs *HTTPServer) Name() string {
	return rs.name
}

func (rs *HTTPServer) InitFlags() {
	prefix := rs.name + "-"
	flag.StringVar(&rs.Host, prefix + "http-host", defaultHost, "HTTP host name to bind with")
	flag.StringVar(&rs.Port, prefix + "http-port", defaultPort, "HTTP port to bind with")
	flag.DurationVar(&rs.ReadTimeout, prefix + "http-read-timeout", defaultReadTimeout, "HTTP read timeout")
	flag.DurationVar(&rs.WriteTimeout, prefix + "http-write-timeout", defaultWriteTimeout, "HTTP write timeout")
	flag.DurationVar(&rs.IdleTimeout, prefix + "http-idle-timeout", defaultIdleTimeout, "HTTP idle timeout")
}

func (rs *HTTPServer) Run() error {
	rs.server = &http.Server{
		Addr: ":" + rs.Port, 
		Handler: rs.router, 
		ReadTimeout: rs.ReadTimeout, 
		WriteTimeout: rs.WriteTimeout, 
		IdleTimeout: rs.IdleTimeout,
	}
	rs.logger.Log.Info("Starting HTTP server on port " + rs.Port)

	if err := rs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed { 
		rs.logger.Log.Error("Error starting HTTP server: ", err)
		return err 
	}

	return nil
}

func (rs *HTTPServer) Stop() <-chan error {
	c := make(chan error, 1)

	go func () {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := rs.server.Shutdown(ctx)
		if err == nil {
			rs.logger.Log.Info("HTTP server stopped.")
		}

		c <- rs.server.Shutdown(ctx)
	}()

	return c
}