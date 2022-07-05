package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	config "github.com/Kaborda-Irina/Kubernetes-Hasher/internal/configs"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/initialize"
	"github.com/joho/godotenv"
)

func main() {
	//Load values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	//Initialize config
	_, logger, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error during loading from config file", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		signal.Stop(sig)
		cancel()
	}()

	initialize.Initialize(ctx, logger, sig)
}
