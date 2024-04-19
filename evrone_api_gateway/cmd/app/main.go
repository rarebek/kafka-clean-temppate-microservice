package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"evrone_api_gateway/internal/app"
	configpkg "evrone_api_gateway/internal/pkg/config"
)

func main() {
	// config
	config, err := configpkg.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// app
	app, err := app.NewApp(*config)
	if err != nil {
		log.Fatal(err)
	}

	// run application
	go func() {
		app.Logger.Info("Listen: ", zap.String("address", config.Server.Host+config.Server.Port))
		if err := app.Run(); err != nil {
			app.Logger.Error("app run", zap.Error(err))
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	// app stops
	app.Logger.Info("api gateway service stops")
	app.Stop()
}
