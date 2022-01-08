package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"star-wars/pkg/server"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("application")
	viper.AddConfigPath("/app")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath("$HOME/config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("couldn't read config file 'application.yaml': %v\n", err.Error())
	}

	viper.AutomaticEnv()
	
	app := server.NewApp()
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		app.Start(ctx)
	}()

	gracefullyShutdown(app, cancel)
}

func gracefullyShutdown(app *server.App, cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	<-signalChan
	log.Print("os.Interrupt - shutting down...\n")

	go func() {
		<-signalChan
		log.Fatal("os.Kill - terminating...\n")
	}()

	gracefulCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := app.Shutdown(gracefulCtx); err != nil {
		log.Printf("shutdown error: %v\n", err)
		defer os.Exit(1)
	} else {
		log.Printf("gracefully stopped\n")
	}
	cancel()
	defer os.Exit(0)

}
