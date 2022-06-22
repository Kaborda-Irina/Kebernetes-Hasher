package initialize

import (
	"context"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/services"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/repositories"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

func Initialize(ctx context.Context, logger *logrus.Logger, sig chan os.Signal, dirPath, algorithm string) {
	// InitializeDB PostgreSQL
	logger.Info("Starting database connection")
	db, err := repositories.InitializeDB(logger)
	if err != nil {
		logger.Error("Failed to connection to db", err)
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

	ticker := time.NewTicker(30 * time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			if service.CheckIsEmptyDB() {
				logger.Info("Empty DB, save data")
				err := service.Start(ctx, dirPath, sig, kuberData)
				if err != nil {
					logger.Fatalf("Error when starting to get hash data %s", err)
				}
			} else {
				logger.Info("Checking, not empty DB")
				for range ticker.C {
					err := service.Check(ctx, ticker, dirPath, sig, kuberData)
					if err != nil {
						logger.Fatalf("Error when starting to check hash data %s", err)
					}
				}
			}

		}
	}(ctx)
	wg.Wait()
}
