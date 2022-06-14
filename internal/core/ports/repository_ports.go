package ports

import (
	"context"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/models"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/pkg/api"
)

//go:generate mockgen -source=repository_ports.go -destination=mocks/mock_repository.go

type IAppRepository interface {
	CheckIsEmptyDB() (bool, error)
}

type IHashRepository interface {
	SaveHashData(ctx context.Context, allHashData []api.HashData) error
	GetHashSum(ctx context.Context, dirFiles string, algorithm string) ([]models.HashDataFromDB, error)
	UpdateDeletedItems(deletedItems []models.DeletedHashes) error
	DeleteAllRowsDB() error
}
