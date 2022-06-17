package main

import (
	"context"
	"flag"
	config "github.com/Kaborda-Irina/Kubernetes-Hasher/internal/configs"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/initialize"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

var dirPath string
var algorithm string

//initializes the binding of the flag to a variable that must run before the main() function
func init() {
	flag.StringVar(&dirPath, "d", "", "a specific file or directory")
	flag.StringVar(&algorithm, "a", "SHA256", "algorithm MD5, SHA1, SHA224, SHA256, SHA384, SHA512, default: SHA256")
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	flag.Parse()

	//Initialize config
	cfg, logger, err := config.LoadConfig()
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

	//dirPath := "../h/h1"
	algorithm := "sha256"
	initialize.Initialize(ctx, cfg, logger, sig, dirPath, algorithm)
}
