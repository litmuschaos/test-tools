package client

import (
	"fmt"
	kubernetes "litmus-helm-agent/pkg/k8s"
	"os"
	"strconv"
	"time"

	types "github.com/litmuschaos/litmusctl/pkg/types"
)

func prepareNewInfra() types.Infra {
	var newInfra types.Infra
	newInfra.Namespace = os.Getenv("NAMESPACE")
	newInfra.SkipSSL, _ = strconv.ParseBool(os.Getenv("SKIP_SSL"))
	newInfra.Mode = os.Getenv("INFRA_MODE")

	// -- OPTIONAL -- //
	newInfra.InfraType = os.Getenv("INFRA_TYPE")
	newInfra.NodeSelector = os.Getenv("INFRA_NODE_SELECTOR")
	newInfra.PlatformName = os.Getenv("PLATFORM_NAME")
	newInfra.ServiceAccount = os.Getenv("SERVICE_ACCOUNT_NAME")
	newInfra.SAExists, _ = strconv.ParseBool(os.Getenv("SA_EXISTS"))
	newInfra.NsExists, _ = strconv.ParseBool(os.Getenv("NS_EXISTS"))
	return newInfra
}

func prepareInfraConfigMap() map[string]string {
	configMapData := make(map[string]string)
	configMapData["SERVER_ADDR"] = os.Getenv("LITMUS_BACKEND_URL")
	configMapData["VERSION"] = os.Getenv("APP_VERSION")
	configMapData["IS_INFRA_CONFIRMED"] = "false"
	configMapData["START_TIME"] = strconv.FormatInt(time.Now().Unix(), 10)
	selector := `["litmuschaos.io/app=chaos-exporter", "litmuschaos.io/app=chaos-operator", "litmuschaos.io/app=event-tracker", "litmuschaos.io/app=workflow-controller"]`
	configMapData["COMPONENTS"] = "DEPLOYMENTS: " + selector
	//TODO fix this line
	configMapData["INFRA_SCOPE"] = os.Getenv("INFRA_MODE")
	configMapData["SKIP_SSL_VERIFY"] = os.Getenv("SKIP_SSL")

	return configMapData
}

func prepareWorkflowControllerConfigMap(clusterID string) map[string]string {
	configMapWorkflowController := make(map[string]string)
	configMapWorkflowController["config"] = (`    containerRuntimeExecutor: ` + os.Getenv("CONTAINER_RUNTIME_EXECUTOR") + `
    executor:
      imagePullPolicy: IfNotPresent
    instanceID: ` + clusterID)
	return configMapWorkflowController

}

func CreateInfra(infraID, accessKey string) {
	clientset := kubernetes.ConnectKubeApi()

	configMap := prepareInfraConfigMap()
	kubernetes.CreateConfigMap(os.Getenv("INFRA_CONFIGMAP_NAME"), configMap, os.Getenv("NAMESPACE"), clientset)

	workflowConfigMap := prepareWorkflowControllerConfigMap(infraID)
	kubernetes.CreateConfigMap(os.Getenv("WORKFLOW_CONTROLER_CONFIGMAP_NAME"), workflowConfigMap, os.Getenv("NAMESPACE"), clientset)

	fmt.Printf("Infra Successfully declared, starting...\n")
}