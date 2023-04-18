package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"github.com/tebrizetayi/messaging-integration-service/internal/api"
	"github.com/tebrizetayi/messaging-integration-service/internal/messagingclients/whatsapp"
)

func main() {
	config := initConfig()

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	messengerClient, _ := whatsapp.NewClient(
		"994552178732",
		config.App.WhatsappAccessToken,
		"https://graph.facebook.com/oauth/access_token",
		"https://graph.facebook.com/v16.0/")

	// Services
	controller := api.NewController(messengerClient)

	// Start the HTTP service listening for requests.
	api := http.Server{
		Addr:           fmt.Sprintf(":%s", config.App.Port),
		Handler:        api.NewAPI(controller),
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Printf("main : API Listening %s", config.App.Port)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown
	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("main : Error starting server: %+v", err)

	case sig := <-shutdown:
		log.Printf("main : %v : Start shutdown..", sig)
	}
}

type Config struct {
	App AppConfig
}
type AppConfig struct {
	Port                string
	WhatsappAccessToken string
	VerifyToken         string
}

func initConfig() Config {
	viper.AutomaticEnv()

	return Config{
		App: AppConfig{
			Port:                viper.GetString("PORT"),
			WhatsappAccessToken: viper.GetString("WHATSAPP_ACCESS_TOKEN"),
			VerifyToken:         viper.GetString("VERIFY_TOKEN"),
		},
	}
}
