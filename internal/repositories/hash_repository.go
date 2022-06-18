package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/consts"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/models"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/pkg/api"

	"github.com/sirupsen/logrus"
)

const nameTable = "hashfiles"

type HashRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewHashRepository(db *sql.DB, logger *logrus.Logger) *HashRepository {
	return &HashRepository{
		db:     db,
		logger: logger,
	}
}

// SaveHashData iterates through all elements of the slice and triggers the save to database function
func (hr HashRepository) SaveHashData(ctx context.Context, allHashData []api.HashData, deploymentData models.DeploymentData) error {
	//_, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	//defer cancel()

	tx, err := hr.db.Begin()
	if err != nil {
		hr.logger.Error("err while saving data in db ", err)
		return err
	}
	query := fmt.Sprintf(`
		INSERT INTO hashfiles (file_name,full_file_path,hash_sum,algorithm,name_pod,name_container,image_tag,time_of_creation) 
		VALUES($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT (full_file_path,algorithm) 
		DO UPDATE SET hash_sum=EXCLUDED.hash_sum`)

	for _, hash := range allHashData {
		_, err = tx.Exec(query, hash.FileName, hash.FullFilePath, hash.Hash, hash.Algorithm, deploymentData.NamePod, deploymentData.NameContainer, deploymentData.Image, deploymentData.Timestamp)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				hr.logger.Error("err in Rollback", err)
				return err
			}
			hr.logger.Error("err while save data in db ", err)
			return err
		}
	}

	return tx.Commit()
}

// GetHashSum retrieves data from the database using the path and algorithm
func (hr HashRepository) GetHashSum(ctx context.Context, dirFiles, algorithm string) ([]models.HashDataFromDB, error) {
	_, cancel := context.WithTimeout(ctx, consts.TimeOut*time.Second)
	defer cancel()

	var allHashDataFromDB []models.HashDataFromDB

	query := fmt.Sprintf("SELECT id,file_name,full_file_path,hash_sum,algorithm FROM %s WHERE full_file_path LIKE $1 and algorithm=$2", nameTable)

	rows, err := hr.db.Query(query, "%"+dirFiles+"%", algorithm)
	if err != nil {
		hr.logger.Error(err)
		return []models.HashDataFromDB{}, err
	}
	for rows.Next() {
		var hashDataFromDB models.HashDataFromDB
		err := rows.Scan(&hashDataFromDB.ID, &hashDataFromDB.FileName, &hashDataFromDB.FullFilePath, &hashDataFromDB.Hash, &hashDataFromDB.Algorithm)
		if err != nil {
			hr.logger.Error(err)
			return []models.HashDataFromDB{}, err
		}
		allHashDataFromDB = append(allHashDataFromDB, hashDataFromDB)
	}

	return allHashDataFromDB, nil
}

// UpdateDeletedItems changes the deleted field to true in the database for each row if the file name has been deleted
func (hr HashRepository) UpdateDeletedItems(deletedItems []models.DeletedHashes) error {
	tx, err := hr.db.Begin()
	if err != nil {
		hr.logger.Error(err)
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET deleted = true WHERE full_file_path=$1 AND algorithm=$2", nameTable)

	for _, item := range deletedItems {
		_, err := tx.Exec(query, item.FilePath, item.Algorithm)

		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			hr.logger.Error(err)
			return err
		}
	}

	return tx.Commit()
}

func (hr HashRepository) DeleteAllRowsDB() error {
	_, err := hr.db.Query("DELETE FROM hashfiles;")
	if err != nil {
		hr.logger.Error("err while deleting rows in db", err)
		return err
	}
	return nil
}
