package initialize

import (
	"context"
	"flag"
	"fmt"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/models"
	"os"

	config "github.com/Kaborda-Irina/Kubernetes-Hasher/internal/configs"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/services"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/repositories"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/pkg/api"

	"github.com/sirupsen/logrus"
)

func Initialize(ctx context.Context, cfg *config.Config, logger *logrus.Logger, sig chan os.Signal, doHelp bool, dirPath, algorithm, checkHashSumFile string) {
	// InitializeDB PostgreSQL
	logger.Info("Starting db connection")
	connectionDB := models.ConnectionDB{
		Dbdriver:   os.Getenv("DB_DRIVER"),
		DbUser:     os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASSWORD"),
		DbPort:     os.Getenv("DB_PORT"),
		DbHost:     os.Getenv("DB_HOST"),
		DbName:     os.Getenv("DB_NAME"),
	}

	db, err := repositories.InitializeDB(connectionDB, logger)
	if err != nil {
		logger.Error("Failed to connection to Postgres", err)
	} else {
		logger.Info("Postgres connection is successful")
	}

	// InitializeDB repository
	repository := repositories.NewAppRepository(db, logger)

	// InitializeDB service
	service, err := services.NewAppService(repository, algorithm, logger)
	if err != nil {
		logger.Fatalf("can't init service: %s", err)
	}

	jobs := make(chan string)
	results := make(chan api.HashData)

	// server
	//newServer := server.NewServer(db)
	//newServer.InitializeDB(cfg)

	switch {
	// InitializeDB custom -h flag
	case doHelp:
		customHelpFlag()
		return
	// InitializeDB custom -d flag
	case len(dirPath) > 0:
		err := service.StartGetHashData(ctx, dirPath, jobs, results, sig)
		if err != nil {
			logger.Error("Error when starting to get hash data ", err)
			return
		}
		return
	// InitializeDB custom -c flag
	case len(checkHashSumFile) > 0:
		err := service.StartCheckHashData(ctx, checkHashSumFile, jobs, results, sig)
		if err != nil {
			logger.Error("Error when starting to check hash data ", err)
			return
		}
		return
	// If the user has not entered a flag
	default:
		logger.Println("use the -h flag on the command line to see all the flags in this app")
	}
}

func customHelpFlag() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Custom help %s:\nYou can use the following flag:\n", os.Args[0])

		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(os.Stderr, "  flag -%v \n       %v\n", f.Name, f.Usage)
		})
	}
	flag.Usage()
}
