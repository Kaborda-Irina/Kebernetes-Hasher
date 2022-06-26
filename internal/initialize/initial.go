package initialize

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/services"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/repositories"
	"github.com/sirupsen/logrus"
)

func Initialize(ctx context.Context, logger *logrus.Logger, sig chan os.Signal, dirPath, algorithm string) {
	// InitializeDB PostgreSQL
	logger.Info("Starting database connection")
	db, err := repositories.InitializeDB(logger)
	if err != nil {
		logger.Error("Failed to connection to database", err)
	}

	// Initialize repository
	repository := repositories.NewAppRepository(db, logger)

	// Initialize service
	service, err := services.NewAppService(repository, algorithm, logger)
	if err != nil {
		logger.Fatalf("can't init service: %s", err)
	}

	// Initialize kubernetesAPI
	kuberData, err := service.ConnectionToKuberAPI()
	if err != nil {
		logger.Fatalf("can't connection to kuberAPI: %s", err)
	}

	duration, err := strconv.Atoi(os.Getenv("DURATION_TIME"))
	if err != nil {
		duration = 15
	}
	ticker := time.NewTicker(time.Duration(duration) * time.Second)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context, ticker *time.Ticker) {
		defer wg.Done()
		for {
			if service.CheckIsEmptyDB(kuberData) {
				logger.Info("Empty DB, save data")
				err := service.Start(ctx, dirPath, sig, kuberData)
				if err != nil {
					logger.Fatalf("Error when starting to get hash data %s", err)
				}
			} else {
				logger.Info("Checking, not empty DB")
				for range ticker.C {
					err := service.Check(ctx, dirPath, sig, kuberData)
					if err != nil {
						logger.Fatalf("Error when starting to check hash data %s", err)
					}
					logger.Info("Check completed")
				}
			}
		}
	}(ctx, ticker)
	wg.Wait()
	ticker.Stop()
}
