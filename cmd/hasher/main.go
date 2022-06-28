package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	config "github.com/Kaborda-Irina/Kubernetes-Hasher/internal/configs"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/initialize"
	"github.com/joho/godotenv"
)

var dirPath string
var algorithm string

//Initializes the binding of the flag to a variable that must run before the main() function
func init() {
	flag.StringVar(&dirPath, "d", "", "a specific file or directory")
	flag.StringVar(&algorithm, "a", "SHA256", "algorithm MD5, SHA1, SHA224, SHA256, SHA384, SHA512, default: SHA256")

	//Load values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	flag.Parse()

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

	initialize.Initialize(ctx, logger, sig, dirPath, algorithm)
}
