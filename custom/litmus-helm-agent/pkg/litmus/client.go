package client

import (
	"fmt"
	"io/ioutil"
	kubernetes "litmus-helm-agent/pkg/k8s"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	ymlToJson "github.com/ghodss/yaml"
	"github.com/golang-jwt/jwt"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/apis/infrastructure"
	types "github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

func prepareNewInfra() types.Infra {
	var newInfra types.Infra
	newInfra.InfraName = os.Getenv("INFRA_NAME")
	newInfra.Namespace = os.Getenv("NAMESPACE")
	newInfra.EnvironmentID = os.Getenv("LITMUS_ENVIRONMENT_ID")
	newInfra.Description = os.Getenv("INFRA_DESCRIPTION")
	newInfra.ProjectId = os.Getenv("LITMUS_PROJECT_ID")
	newInfra.Mode = os.Getenv("INFRA_MODE")
	newInfra.SkipSSL, _ = strconv.ParseBool(os.Getenv("SKIP_SSL"))

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

func prepareInfraSecret(infraConnect infrastructure.RegisterInfra, accessKey string) map[string][]byte {
	secretData := make(map[string][]byte)
	InfraID := infraConnect.RegisterInfraDetails.InfraID
	secretData["INFRA_ID"] = []byte(InfraID)
	secretData["ACCESS_KEY"] = []byte(accessKey)

	return secretData
}

func prepareWorkflowControllerConfigMap(clusterID string) map[string]string {
	configMapWorkflowController := make(map[string]string)
	configMapWorkflowController["config"] = (`    containerRuntimeExecutor: ` + os.Getenv("CONTAINER_RUNTIME_EXECUTOR") + `
    executor:
      imagePullPolicy: IfNotPresent
    instanceID: ` + clusterID)
	return configMapWorkflowController

}

func GetProjectID(credentials types.Credentials) string {
	var result string
	userDetails, err := apis.GetProjectDetails(credentials)
	if err != nil {
		fmt.Printf("Error, cannot get project details: " + err.Error())
		os.Exit(1)
	}

	for _, project := range userDetails.Data.Projects {
		for _, member := range project.Members {
			if (member.UserID == userDetails.Data.ID) && (member.Role == "Owner" || member.Role == "Editor") {
				result = project.ID
				break
			}
		}
	}
	if result == "" {
		utils.Red.Println("\n❌ No project found with owner or editor access to current user" + "\n")
		os.Exit(1)
	}
	return result
}

func GetInfraWithName(credentials types.Credentials, searchInfra types.Infra) (models.Infra, error) {
	infras, err := infrastructure.GetInfraList(credentials, searchInfra.ProjectId, models.ListInfraRequest{})
	if err != nil {
		return models.Infra{}, err
	}
	for _, infra := range infras.Data.ListInfraDetails.Infras {
		if infra.Name == searchInfra.InfraName {
			return *infra, nil
		}
	}
	return models.Infra{}, nil
}

func CreateInfra(credentials types.Credentials) {
	newInfra := prepareNewInfra()

	if newInfra.ProjectId == "" {
		newInfra.ProjectId = GetProjectID(credentials)
	}

	infraExist, err := GetInfraWithName(credentials, newInfra)
	if err != nil {
		utils.Red.Printf("\n❌ Error, cannot search if infrastructure exist: %v", err.Error())
		os.Exit(1)
	}

	if reflect.ValueOf(infraExist).IsZero() {
		connectionData, err := infrastructure.ConnectInfra(newInfra, credentials)

		if err != nil {
			utils.Red.Println("\n❌ Infrastructure registration failed: " + err.Error() + "\n")
			os.Exit(1)
		}
		if (connectionData.Data == infrastructure.RegisterInfra{}) {
			fmt.Printf("❌ Agent empty: Registration failed did graphql change ? \n")
			os.Exit(1)
		}

		if connectionData.Data.RegisterInfraDetails.Token == "" {
			utils.Red.Println("\n❌ failed to get the infrastructure registration token: " + "\n")
			os.Exit(1)
		}

		accessKey, err := validateInfra(connectionData.Data.RegisterInfraDetails.Token, credentials.Endpoint)

		if err != nil {
			utils.Red.Println("❌ Error validate infrastructure: ", err)
			os.Exit(1)
		}

		clientset := kubernetes.ConnectKubeApi()
		configMap := prepareInfraConfigMap()
		kubernetes.CreateConfigMap(os.Getenv("INFRA_CONFIGMAP_NAME"), configMap, os.Getenv("NAMESPACE"), clientset)

		secret := prepareInfraSecret(connectionData.Data, accessKey)
		kubernetes.CreateSecret(os.Getenv("INFRA_SECRET_NAME"), secret, os.Getenv("NAMESPACE"), clientset)

		workflowConfigMap := prepareWorkflowControllerConfigMap(connectionData.Data.RegisterInfraDetails.InfraID)
		kubernetes.CreateConfigMap(os.Getenv("WORKFLOW_CONTROLER_CONFIGMAP_NAME"), workflowConfigMap, os.Getenv("NAMESPACE"), clientset)

		fmt.Printf("Infra Successfully declared, starting...\n")
	} else {
		fmt.Printf("Infra already exist, starting...\n")
	}
}

func validateInfra(token string, endpoint string) (string, error) {
	var accessKey string

	path := fmt.Sprintf("%s/%s/%s.yaml", endpoint, utils.ChaosYamlPath, token)
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return accessKey, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return accessKey, err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return accessKey, err
	}
	manifests := strings.Split(string(resp_body), "---")
	for _, manifest := range manifests {
		if len(strings.TrimSpace(manifest)) > 0 {
			jsonValue, err := ymlToJson.YAMLToJSON([]byte(manifest))
			if err != nil {
				return accessKey, err
			}
			fieldName, _, _, err := jsonparser.Get([]byte(jsonValue), "metadata", "name")
			if err != nil {
				return accessKey, err
			}
			fieldKind, _, _, err := jsonparser.Get([]byte(jsonValue), "kind")
			if err != nil {
				return accessKey, err
			}
			if string(fieldName) == "subscriber-secret" && string(fieldKind) == "Secret" {
				if fieldData, _, _, err := jsonparser.Get([]byte(jsonValue), "stringData", "ACCESS_KEY"); err != nil {
					return accessKey, err
				} else {
					accessKey = string(fieldData)
				}
			}
		}
	}
	return accessKey, err
}

func Login(LITMUS_FRONTEND_URL string, LITMUS_USERNAME string, LITMUS_PASSWORD string) types.Credentials {
	msg := ""

	if len(LITMUS_FRONTEND_URL) == 0 {
		msg = msg + "LITMUS_FRONTEND_URL, "
	}

	if len(LITMUS_USERNAME) == 0 {
		msg = msg + "LITMUS_USERNAME, "
	}

	if len(LITMUS_PASSWORD) == 0 {
		msg = msg + "LITMUS_PASSWORD, "
	}
	if msg != "" {
		utils.Red.Println("❌ " + msg + " should be set as env var")
		os.Exit(1)
	}

	var authInput types.AuthInput
	authInput.Endpoint = LITMUS_FRONTEND_URL
	authInput.Username = LITMUS_USERNAME
	authInput.Password = LITMUS_PASSWORD

	resp, err := apis.Auth(authInput)
	utils.PrintError(err)
	// Decoding token
	token, _ := jwt.Parse(resp.AccessToken, nil)
	if token == nil {
		utils.Red.Println("\n❌ Cannot get token for user: " + authInput.Username + "\n")
		os.Exit(1)
	}

	var credentials types.Credentials
	credentials.Username = authInput.Username
	credentials.Endpoint = authInput.Endpoint
	credentials.ServerEndpoint = authInput.Endpoint
	credentials.Token = resp.AccessToken

	return credentials
}
