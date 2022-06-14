package ports

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/models"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/pkg/api"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=service_ports.go -destination=mocks/mock_service.go

type IAppService interface {
	Start(ctx context.Context, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) error
	Check(ctx context.Context, ticker *time.Ticker, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) error
}

type IHashService interface {
	SaveHashData(ctx context.Context, allHashData []api.HashData) error
	GetHashSum(ctx context.Context, dirFiles string) ([]models.HashDataFromDB, error)
	ChangedHashes(ctx context.Context, ticker *time.Ticker, currentHashData []api.HashData, hashSumFromDB []models.HashDataFromDB) error
	CreateHash(path string) api.HashData
	WorkerPool(ctx context.Context, countWorkers int, jobs chan string, results chan api.HashData, logger *logrus.Logger)
	Worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- api.HashData, logger *logrus.Logger)
}
