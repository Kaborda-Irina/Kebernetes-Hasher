package repositories

import (
	"database/sql"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/ports"
	"github.com/sirupsen/logrus"
)

type AppRepository struct {
	ports.IHashRepository
	db     *sql.DB
	logger *logrus.Logger
}

func NewAppRepository(db *sql.DB, logger *logrus.Logger) *AppRepository {
	return &AppRepository{
		IHashRepository: NewHashRepository(db, logger),
		db:              db,
		logger:          logger,
	}
}

func (ar AppRepository) CheckIsEmptyDB() (bool, error) {
	var count int
	row := ar.db.QueryRow("SELECT COUNT(*) FROM hashfiles LIMIT 1")
	err := row.Scan(&count)
	if err != nil {
		ar.logger.Info("err while scan row in db ", err)
		return false, err
	}

	if count < 1 {
		return true, nil
	}
	return false, nil
}
