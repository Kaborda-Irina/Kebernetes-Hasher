package services

import (
	"context"
	"time"

	"os"
	"strings"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/ports"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/repositories"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/pkg/api"

	"github.com/sirupsen/logrus"
)

type AppService struct {
	ports.IHashService
	ports.IAppRepository
	logger *logrus.Logger
}

// NewAppService creates a new struct AppService
func NewAppService(r *repositories.AppRepository, algorithm string, logger *logrus.Logger) (*AppService, error) {
	algorithm = strings.ToUpper(algorithm)
	IHashService, err := NewHashService(r.IHashRepository, algorithm, logger)
	if err != nil {
		return nil, err
	}
	return &AppService{
		IHashService:   IHashService,
		IAppRepository: r,
		logger:         logger,
	}, nil
}
func (as *AppService) LaunchHasher(ctx context.Context, flagName string, sig chan os.Signal) []api.HashData {
	jobs := make(chan string)
	results := make(chan api.HashData)
	go as.IHashService.WorkerPool(ctx, jobs, results, as.logger)
	go api.SearchFilePath(ctx, flagName, jobs, as.logger)
	allHashData := api.Result(ctx, results, sig)

	return allHashData
}

func (as *AppService) CheckIsEmptyDB() bool {
	isEmptyDB, err := as.IAppRepository.CheckIsEmptyDB()
	if err != nil {
		as.logger.Fatalf("error while saving data to db %s", err)
	}
	return isEmptyDB
}

// StartGetHashData getting the hash sum of all files, outputs to os.Stdout and saves to the database
func (as *AppService) Start(ctx context.Context, flagName string, sig chan os.Signal) error {
	allHashData := as.LaunchHasher(ctx, flagName, sig)
	err := as.IHashService.SaveHashData(ctx, allHashData)
	if err != nil {
		as.logger.Error("Error save hash data to database ", err)
		return err
	}
	return nil
}

// StartCheckHashData getting the hash sum of all files, matches them and outputs to os.Stdout changes
func (as *AppService) Check(ctx context.Context, ticker *time.Ticker, flagName string, sig chan os.Signal) error {
	allHashDataCurrent := as.LaunchHasher(ctx, flagName, sig)
	allHashDataFromDB, err := as.IHashService.GetHashSum(ctx, flagName)
	if err != nil {
		as.logger.Error("Error getting hash data from database ", err)
		return err
	}
	err = as.IHashService.ChangedHashes(ctx, ticker, allHashDataCurrent, allHashDataFromDB)
	if err != nil {
		as.logger.Error("Error match data currently and data from db ", err)
		return err
	}
	return nil
}
