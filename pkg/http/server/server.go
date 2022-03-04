package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Kolakanmi/grey_wallet/pkg/envconfig"
)

type (
	//Config - struct
	Config struct {
		Address           string        `envconfig:"HTTP_ADDRESS"`
		Port              string        `envconfig:"HTTP_PORT" default:"8080"`
		ReadTimeout       time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"10s"`
		ReadHeaderTimeout time.Duration `envconfig:"HTTP_READ_HEADER_TIMEOUT" default:"20s"`
		WriteTimeout      time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"120s"`
		IdleTimeout       time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" default:"180s"`
		ShutdownTimeout   time.Duration `envconfig:"HTTP_SHUTDOWN_TIMEOUT" default:"10s"`
	}
)

//ListenAndServe - listen function
func ListenAndServe(conf Config, handler http.Handler) {
	port := conf.Port
	if port == "" {
		appEnginePort := os.Getenv("PORT")
		if appEnginePort == "" {
			port = "80"
		}
	}
	address := fmt.Sprintf("%s:%s", conf.Address, port)

	srv := http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	}

	log.Printf("HTTP Server is listening on: %s", address)
	//Run server in goroutine to implement graceful shutdown.
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("listen: %v\n", err)
		}
	}()
	//Graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	//Block until signal is received
	<-signals
	ctx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
	defer cancel()
	log.Print("Server shutting down!!!")
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Panicf("http server shutdown with error: %v", err)
	}
}

//LoadConfigFromEnv - load from env
func LoadConfigFromEnv() *Config {
	var config Config
	envconfig.Load(&config)
	return &config
}
