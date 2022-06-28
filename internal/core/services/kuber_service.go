package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Kaborda-Irina/Kubernetes-Hasher/internal/core/models"
	"github.com/sirupsen/logrus"
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

func (ks *KuberService) ConnectionToKuberAPI() (models.KuberData, error) {
	// Connect to Kubernetes API
	ks.logger.Info("### 🌀 Attempting to use in cluster config")
	config, err := rest.InClusterConfig()
	if err != nil {
		ks.logger.Error(err)
		return models.KuberData{}, err
	}

	ks.logger.Info("### 💻 Connecting to Kubernetes API, using host: ", config.Host)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		ks.logger.Error(err)
		return models.KuberData{}, err
	}

	namespaceBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		ks.logger.Error(err)
		return models.KuberData{}, err
	}
	namespace := string(namespaceBytes)

	podName := os.Getenv("POD_NAME")

	deploymentName := func(podName string) string {
		elements := strings.Split(podName, "-")
		newElements := elements[:len(elements)-2]
		return strings.Join(newElements, "-")
	}(podName)
	if deploymentName == "" {
		ks.logger.Fatalln("### 💥 Env var DEPLOYMENT_NAME was not set")
	}
	deploymentType := os.Getenv("DEPLOYMENT_TYPE")
	kuberData := models.KuberData{
		Clientset:  clientset,
		Namespace:  namespace,
		TargetName: deploymentName,
		TargetType: deploymentType,
	}
	return kuberData, nil
}

func (ks *KuberService) GetDataFromKuberAPI(kuberData models.KuberData) (models.DeploymentData, error) {
	allDeploymentData, err := kuberData.Clientset.AppsV1().Deployments(kuberData.Namespace).Get(context.Background(), kuberData.TargetName, metav1.GetOptions{})
	if err != nil {
		ks.logger.Error("err while getting data from kuberAPI ", err)
		return models.DeploymentData{}, err
	}

	deploymentData := models.DeploymentData{}
	deploymentData.NameDeployment = kuberData.TargetName
	deploymentData.Timestamp = fmt.Sprintf("%v", allDeploymentData.CreationTimestamp)
	deploymentData.NamePod = os.Getenv("POD_NAME")

	for _, v := range allDeploymentData.Spec.Template.Spec.Containers {
		deploymentData.Image = v.Image
	}

	return deploymentData, nil
}

func (ks *KuberService) RolloutDeployment(kuberData models.KuberData) error {
	patchData := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}}}`, time.Now().Format(time.RFC3339))
	_, err := kuberData.Clientset.AppsV1().Deployments(kuberData.Namespace).Patch(context.Background(), kuberData.TargetName, types.StrategicMergePatchType, []byte(patchData), metav1.PatchOptions{FieldManager: "kubectl-rollout"})

	if err != nil {
		ks.logger.Printf("### 👎 Warning: Failed to patch %v, restart failed: %v", kuberData.TargetType, err)
		return err
	} else {
		ks.logger.Printf("### ✅ Target %v, named %v was restarted!", kuberData.TargetType, kuberData.TargetName)
	}
	return nil
}
