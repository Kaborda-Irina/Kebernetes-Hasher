package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/consts"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/models"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/ports"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/pkg/api"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/pkg/hasher"

	"github.com/sirupsen/logrus"
)

type HashService struct {
	hashRepository ports.IHashRepository
	hasher         hasher.IHasher
	alg            string
	logger         *logrus.Logger
}

// NewHashService creates a new struct HashService
func NewHashService(hashRepository ports.IHashRepository, alg string, logger *logrus.Logger) (*HashService, error) {
	h, err := hasher.NewHashSum(alg)
	if err != nil {
		return nil, err
	}
	return &HashService{
		hashRepository: hashRepository,
		hasher:         h,
		alg:            alg,
		logger:         logger,
	}, nil
}

// WorkerPool launches a certain number of workers for concurrent processing
func (hs HashService) WorkerPool(ctx context.Context, jobs chan string, results chan api.HashData, logger *logrus.Logger) {
	ctx, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	for w := 1; w <= consts.CountWorkers; w++ {
		wg.Add(1)
		go hs.Worker(ctx, &wg, jobs, results, logger)
	}
	defer close(results)
	wg.Wait()
}

// Worker gets jobs from a pipe and writes the result to stdout and database
func (hs HashService) Worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- api.HashData, _ *logrus.Logger) {
	_, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	defer cancel()
	defer wg.Done()
	for j := range jobs {
		results <- hs.CreateHash(j)
	}
}

// CreateHash creates a new object with a hash sum
func (hs HashService) CreateHash(path string) api.HashData {
	file, err := os.Open(path)
	if err != nil {
		hs.logger.Error("not exist file path", err)
		return api.HashData{}
	}
	defer file.Close()

	outputHashSum := api.HashData{}
	res, err := hs.hasher.Hash(file)
	if err != nil {
		hs.logger.Error("not got hash sum", err)
		return api.HashData{}
	}
	outputHashSum.Hash = res
	outputHashSum.FileName = filepath.Base(path)
	outputHashSum.FullFilePath = path
	outputHashSum.Algorithm = hs.alg
	return outputHashSum
}

// SaveHashData accesses the repository to save data to the database
func (hs HashService) SaveHashData(ctx context.Context, allHashData []api.HashData, deploymentData models.DeploymentData) error {
	ctx, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	defer cancel()

	err := hs.hashRepository.SaveHashData(ctx, allHashData, deploymentData)
	if err != nil {
		hs.logger.Error("error while saving data to db", err)
		return err
	}
	return nil
}

// GetHashData accesses the repository to get data from the database
func (hs HashService) GetHashData(ctx context.Context, dirFiles string) ([]models.HashDataFromDB, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	defer cancel()

	hash, err := hs.hashRepository.GetHashData(ctx, dirFiles, hs.alg)
	if err != nil {
		hs.logger.Error("hash service didn't get hash sum", err)
		return nil, err
	}

	return hash, nil
}

func (hs HashService) TruncateTable() error {
	err := hs.hashRepository.TruncateTable()
	if err != nil {
		hs.logger.Error("err while deleting rows in db", err)
		return err
	}
	return nil
}

// IsDataChanged checks if the current data has changed with the data stored in the database
func (hs HashService) IsDataChanged(ticker *time.Ticker, currentHashData []api.HashData, hashDataFromDB []models.HashDataFromDB, deploymentData models.DeploymentData) (bool, error) {
	isDataChanged := matchwithDataDB(hashDataFromDB, currentHashData, ticker, deploymentData)
	isAddedFiles := matchWithDataCurrent(currentHashData, hashDataFromDB, ticker)

	if isDataChanged || isAddedFiles {
		return true, nil
	}
	return false, nil
}

func matchwithDataDB(hashSumFromDB []models.HashDataFromDB, currentHashData []api.HashData, ticker *time.Ticker, deploymentData models.DeploymentData) bool {
	for _, dataFromDB := range hashSumFromDB {
		trigger := false
		for _, dataCurrent := range currentHashData {
			if dataFromDB.FullFilePath == dataCurrent.FullFilePath && dataFromDB.Algorithm == dataCurrent.Algorithm {
				if dataFromDB.Hash != dataCurrent.Hash {
					fmt.Printf("Changed: file - %s the path %s, old hash sum %s, new hash sum %s\n",
						dataFromDB.FileName, dataFromDB.FullFilePath, dataFromDB.Hash, dataCurrent.Hash)
					ticker.Stop()
					return true
				}
				if dataFromDB.ImageContainer != deploymentData.Image {
					fmt.Printf("Changed image container: file - %s the path %s, old image %s, new image %s\n",
						dataFromDB.FileName, dataFromDB.FullFilePath, dataFromDB.ImageContainer, deploymentData.Image)
					ticker.Stop()
					return true
				}
				trigger = true
				break
			}
		}

		if !trigger {
			fmt.Printf("Deleted: file - %s the path %s hash sum %s\n", dataFromDB.FileName, dataFromDB.FullFilePath, dataFromDB.Hash)
			ticker.Stop()
			return true
		}
	}
	return false
}

func matchWithDataCurrent(currentHashData []api.HashData, hashDataFromDB []models.HashDataFromDB, ticker *time.Ticker) bool {
	dataFromDB := make(map[string]struct{}, len(hashDataFromDB))
	for _, value := range hashDataFromDB {
		dataFromDB[value.FullFilePath] = struct{}{}
	}

	for _, dataCurrent := range currentHashData {
		if _, ok := dataFromDB[dataCurrent.FullFilePath]; !ok {
			fmt.Printf("Added: file - %s the path %s hash sum %s\n",
				dataCurrent.FileName, dataCurrent.FullFilePath, dataCurrent.Hash)
			ticker.Stop()
			return true
		}
	}
	return false
}
