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
	CheckIsEmptyDB() bool
	Start(ctx context.Context, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) error
	Check(ctx context.Context, ticker *time.Ticker, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) error
	RolloutDeployment()
}

type IHashService interface {
	SaveHashData(ctx context.Context, allHashData []api.HashData, deploymentData models.DeploymentData) error
	GetHashSum(ctx context.Context, dirFiles string) ([]models.HashDataFromDB, error)
	TruncateTable() error
	IsDataChanged(ctx context.Context, ticker *time.Ticker, currentHashData []api.HashData, hashSumFromDB []models.HashDataFromDB) (bool, error)
	CreateHash(path string) api.HashData
	WorkerPool(ctx context.Context, jobs chan string, results chan api.HashData, logger *logrus.Logger)
	Worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- api.HashData, logger *logrus.Logger)
}

type IKuberService interface {
	ConnectionToKuberAPI() (models.KuberData, error)
	GetDataFromKuberAPI(kuberData models.KuberData) (models.DeploymentData, error)
	RolloutDeployment(kuberData models.KuberData) error
}
