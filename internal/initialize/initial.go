package initialize

import (
	"context"
	"database/sql"
	"fmt"
	config "github.com/Kaborda-Irina/Kubernetes-Hasher/internal/configs"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/models"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/services"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/repositories"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sync"
	"time"
)

func Initialize(ctx context.Context, cfg *config.Config, logger *logrus.Logger, sig chan os.Signal, dirPath, algorithm string) {
	// InitializeDB PostgreSQL
	logger.Info("Starting database connection")
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
		logger.Error("Failed to connection to db", err)
	}

	// InitializeDB repository
	repository := repositories.NewAppRepository(db, logger)

	// InitializeDB service
	service, err := services.NewAppService(repository, algorithm, logger)
	if err != nil {
		logger.Fatalf("can't init service: %s", err)
	}

	fmt.Println(IsEmptyDB(db))
	ticker := time.NewTicker(5 * time.Second)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			if IsEmptyDB(db) {
				err := service.Start(ctx, dirPath, sig)
				if err != nil {
					logger.Error("Error when starting to get hash data ", err)
					return
				}
			} else {
				fmt.Println("not empty")
				for range ticker.C {
					err := service.Check(ctx, ticker, dirPath, sig)
					if err != nil {
						logger.Error("Error when starting to check hash data ", err)
						os.Exit(1)
					}
				}
			}

		}
	}(ctx)

	wg.Wait()
	fmt.Println("Ticker stopped")
}

func IsEmptyDB(db *sql.DB) bool {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM hashfiles LIMIT 1")
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count < 1 {
		return true
	}
	return false
}
