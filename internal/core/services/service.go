package services

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/models"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/ports"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/repositories"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/pkg/api"

	"github.com/sirupsen/logrus"
)

type AppService struct {
	ports.IHashService
	ports.IAppRepository
	ports.IKuberService
	logger *logrus.Logger
}

// NewAppService creates a new struct AppService
func NewAppService(r *repositories.AppRepository, algorithm string, logger *logrus.Logger) (*AppService, error) {
	algorithm = strings.ToUpper(algorithm)
	IHashService, err := NewHashService(r.IHashRepository, algorithm, logger)
	if err != nil {
		return nil, err
	}
	kuberService := NewKuberService(logger)
	return &AppService{
		IHashService:   IHashService,
		IAppRepository: r,
		IKuberService:  kuberService,
		logger:         logger,
	}, nil
}

//GetPID getting pid by process name
func (as *AppService) GetPID(configData models.ConfigMapData) (int, error) {
	if os.Chdir(os.Getenv("PROC_DIR")) != nil {
		fmt.Println("/proc unavailable.")
		os.Exit(1)
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		fmt.Println("unable to read /proc directory.")
	}
	var pid int
	for _, file := range files {
		if !file.IsDir() {
			return 0, err
		}

		// Our directory name should convert to integer if it's a PID
		pid, err = strconv.Atoi(file.Name())
		if err != nil {
			return 0, err
		}

		// Open the /proc/xxx/stat file to read the name
		f, err := os.Open(file.Name() + "/stat")
		if err != nil {
			fmt.Println("unable to open", file.Name())
		}
		defer f.Close()

		r := bufio.NewReader(f)
		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), configData.ProcName) {
				return pid, nil
			}
		}
	}

	return pid, nil
}

//LaunchHasher takes a path to a directory and returns HashData
func (as *AppService) LaunchHasher(ctx context.Context, dirPath string, sig chan os.Signal) []api.HashData {
	jobs := make(chan string)
	results := make(chan api.HashData)
	go as.IHashService.WorkerPool(ctx, jobs, results, as.logger)
	go api.SearchFilePath(ctx, dirPath, jobs, as.logger)
	allHashData := api.Result(ctx, results, sig)

	return allHashData
}

//IsExistDeploymentNameInDB checks if the database is empty
func (as *AppService) IsExistDeploymentNameInDB(deploymentName string) bool {
	isEmptyDB, err := as.IAppRepository.IsExistDeploymentNameInDB(deploymentName)
	if err != nil {
		as.logger.Fatalf("database check error %s", err)
	}
	return isEmptyDB
}

// StartGetHashData getting the hash sum of all files, outputs to os.Stdout and saves to the database
func (as *AppService) Start(ctx context.Context, dirPath string, sig chan os.Signal, deploymentData models.DeploymentData) error {
	allHashData := as.LaunchHasher(ctx, dirPath, sig)
	err := as.IHashService.SaveHashData(ctx, allHashData, deploymentData)
	if err != nil {
		as.logger.Error("Error save hash data to database ", err)
		return err
	}

	return nil
}

// StartCheckHashData getting the hash sum of all files, matches them and outputs to os.Stdout changes
func (as *AppService) Check(ctx context.Context, dirPath string, sig chan os.Signal, deploymentData models.DeploymentData, kuberData models.KuberData) error {
	hashDataCurrentByDirPath := as.LaunchHasher(ctx, dirPath, sig)

	dataFromDBbyPodName, err := as.IHashService.GetHashData(ctx, dirPath, deploymentData)
	if err != nil {
		as.logger.Error("Error getting hash data from database ", err)
		return err
	}

	isDataChanged, err := as.IHashService.IsDataChanged(hashDataCurrentByDirPath, dataFromDBbyPodName, deploymentData)
	if err != nil {
		as.logger.Error("Error match data currently and data from database ", err)
		return err
	}
	if isDataChanged {
		err := as.IHashService.DeleteFromTable(deploymentData.NameDeployment)
		if err != nil {
			as.logger.Error("Error while deleting rows in database", err)
			return err
		}

		err = as.IKuberService.RolloutDeployment(kuberData)
	}
	return nil
}
