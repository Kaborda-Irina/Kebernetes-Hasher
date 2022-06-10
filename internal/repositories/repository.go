package repositories

import (
	"database/sql"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/ports"

	"github.com/sirupsen/logrus"
)

type AppRepository struct {
	ports.IHashRepository
	logger *logrus.Logger
}

func NewAppRepository(db *sql.DB, logger *logrus.Logger) *AppRepository {
	return &AppRepository{
		IHashRepository: NewHashRepository(db, logger),
		logger:          logger,
	}
}
