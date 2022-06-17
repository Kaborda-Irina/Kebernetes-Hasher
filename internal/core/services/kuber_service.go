package services

import (
	"context"
	"fmt"
	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/models"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KuberService struct {
	logger *logrus.Logger
}

// NewHashService creates a new struct HashService
func NewKuberService(logger *logrus.Logger) *KuberService {
	return &KuberService{
		logger: logger,
	}
}

func (ks *KuberService) ConnectionToKuberAPI() models.KuberData {
	targetName := os.Getenv("DEPLOYMENT_NAME")
	if targetName == "" {
		ks.logger.Fatalln("### 💥 Env var DEPLOYMENT_NAME was not set")
	}
	targetType := os.Getenv("DEPLOYMENT_TYPE")

	// Connect to Kubernetes API
	ks.logger.Info("### 🌀 Attempting to use in cluster config")
	config, err := rest.InClusterConfig()
	if err != nil {
		ks.logger.Fatalln(err)
	}

	ks.logger.Info("### 💻 Connecting to Kubernetes API, using host: %s", config.Host)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		ks.logger.Fatal(err)
	}

	namespaceBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		ks.logger.Fatal(err)
	}
	namespace := string(namespaceBytes)

	kuberData := models.KuberData{
		clientset,
		namespace,
		targetName,
		targetType,
	}
	return kuberData
}

func (ks *KuberService) GetDataFromKuberAPI(kuberData models.KuberData) {
	res, err := kuberData.Clientset.AppsV1().Deployments(kuberData.Namespace).Get(context.Background(), kuberData.TargetName, metav1.GetOptions{})
	if err != nil {
		ks.logger.Info("hhhhhhhhhhhh ", err)
	}
	ks.logger.Info("result get deploy ", res)
}

func (ks *KuberService) RolloutDeployment(kuberData models.KuberData) error {
	patchData := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))
	_, err := kuberData.Clientset.AppsV1().Deployments(kuberData.Namespace).Patch(context.Background(), kuberData.TargetName, types.StrategicMergePatchType, []byte(patchData), metav1.PatchOptions{FieldManager: "kubectl-rollout"})

	if err != nil {
		ks.logger.Info("### 👎 Warning: Failed to patch %s, restart failed: %v", kuberData.TargetType, err)
		return err
	} else {
		ks.logger.Info("### ✅ Target %s, named %s was restarted!", kuberData.TargetType, kuberData.TargetName)
	}
	return nil
}
