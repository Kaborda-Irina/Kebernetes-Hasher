package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
func (hs HashService) SaveHashData(ctx context.Context, allHashData []api.HashData) error {
	ctx, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	defer cancel()

	err := hs.hashRepository.SaveHashData(ctx, allHashData)
	if err != nil {
		hs.logger.Error("error while saving data to db", err)
		return err
	}
	return nil
}

// GetHashSum accesses the repository to get data from the database
func (hs HashService) GetHashSum(ctx context.Context, dirFiles string) ([]models.HashDataFromDB, error) {
	ctx, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	defer cancel()

	hash, err := hs.hashRepository.GetHashSum(ctx, dirFiles, hs.alg)
	if err != nil {
		hs.logger.Error("hash service didn't get hash sum", err)
		return nil, err
	}

	return hash, nil
}

// ChangedHashes checks if the current data has changed with the data stored in the database
func (hs HashService) ChangedHashes(ctx context.Context, ticker *time.Ticker, currentHashData []api.HashData, hashSumFromDB []models.HashDataFromDB) error {
	_, cancel := context.WithCancel(ctx)

	var deletedResult []models.DeletedHashes
	var trigger bool
	var count int

	for _, dataFromDB := range hashSumFromDB {
		trigger = false
		for _, dataCurrent := range currentHashData {
			if dataFromDB.FullFilePath == dataCurrent.FullFilePath {
				if dataFromDB.Hash != dataCurrent.Hash {
					count++
					fmt.Printf("Changed: file - %s the path %s, old hash sum %s, new hash sum %s\n",
						dataFromDB.FileName, dataFromDB.FullFilePath, dataFromDB.Hash, dataCurrent.Hash)
					ticker.Stop()
					cancel()
					os.Exit(1)

				}
				trigger = true
				break
			}
		}

		if !trigger {
			count++
			fmt.Printf("Deleted: file - %s the path %s hash sum %s\n",
				dataFromDB.FileName, dataFromDB.FullFilePath, dataFromDB.Hash)
			deletedResult = append(deletedResult, models.DeletedHashes{
				FileName:    dataFromDB.FileName,
				FilePath:    dataFromDB.FullFilePath,
				OldChecksum: dataFromDB.Hash,
				Algorithm:   dataFromDB.Algorithm,
			})
			ticker.Stop()
			cancel()

		}
	}

	for _, dataCurrent := range currentHashData {
		trigger = false
		for _, dataFromDB := range hashSumFromDB {
			if dataCurrent.FullFilePath == dataFromDB.FullFilePath {
				trigger = true
				break
			}
		}

		if !trigger {
			count++
			fmt.Printf("Added: file - %s the path %s hash sum %s\n",
				dataCurrent.FileName, dataCurrent.FullFilePath, dataCurrent.Hash)
		}
	}

	if len(deletedResult) > 0 {
		err := hs.hashRepository.UpdateDeletedItems(deletedResult)
		if err != nil {
			hs.logger.Error(err)
			return err
		}
	}

	if count == 0 {
		fmt.Println("Files have not been changed, added or removed")
	}
	return nil
}
