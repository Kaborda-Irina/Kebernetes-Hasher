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
	GetPID(confirData models.ConfigMapData) (int, error)
	CheckIsEmptyDB(kuberData models.KuberData) bool
	Start(ctx context.Context, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) error
	Check(ctx context.Context, ticker *time.Ticker, flagName string, jobs chan string, results chan api.HashData, sig chan os.Signal) error
	RolloutDeployment()
}

type IHashService interface {
	SaveHashData(ctx context.Context, allHashData []api.HashData, deploymentData models.DeploymentData) error
	GetHashData(ctx context.Context, dirPath string, deploymentData models.DeploymentData) ([]models.HashDataFromDB, error)
	DeleteFromTable(nameDeployment string) error
	IsDataChanged(currentHashData []api.HashData, hashSumFromDB []models.HashDataFromDB, deploymentData models.DeploymentData) (bool, error)
	CreateHash(path string) api.HashData
	WorkerPool(ctx context.Context, jobs chan string, results chan api.HashData, logger *logrus.Logger)
	Worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- api.HashData, logger *logrus.Logger)
}

type IKuberService interface {
	GetDataFromK8sAPI() (models.KuberData, models.DeploymentData, models.ConfigMapData, error)
	ConnectionToK8sAPI() (models.KuberData, error)
	GetDataFromDeployment(kuberData models.KuberData) (models.DeploymentData, error)
	GetDataFromConfigMap(kuberData models.KuberData, label string) (models.ConfigMapData, error)
	RolloutDeployment(kuberData models.KuberData) error
}
