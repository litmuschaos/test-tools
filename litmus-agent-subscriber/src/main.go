package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"time"
        "strconv"

	"github.com/golang-jwt/jwt"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	types "github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	corev1r "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1r "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//"k8s.io/client-go/kubernetes"
)

var (
	LITMUS_FRONTEND_URL string
	LITMUS_BACKEND_URL  string
	LITMUS_USERNAME     string
	LITMUS_PASSWORD     string
	LITMUS_PROJECT_ID   string

	AGENT_NAME        string
	AGENT_DESCRIPTION string

	NAMESPACE            string
	APP_VERSION          string
	SERVICE_ACCOUNT_NAME string
	CONFIG_MAP_NAME      string
	RELEASE_NAME         string
	ACTION               string

	CLUSTER_ID	     string
)

type agentData struct {
	ClusterID   string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	AccessKey   string `json:"access_key"`
        IsClusterConfirmed bool `json:"is_cluster_confirmed"`
}

type agentResponse struct {
	Data struct {
		GetCluster []agentData `json:"getCluster"`
	} `json:"data"`
}

func init() {
	flag.StringVar(&ACTION, "action", "", "create|delete litmus agent")
	flag.Parse()
	LITMUS_FRONTEND_URL = os.Getenv("LITMUS_FRONTEND_URL")
	LITMUS_BACKEND_URL = os.Getenv("LITMUS_BACKEND_URL")

	LITMUS_USERNAME = os.Getenv("LITMUS_USERNAME")
	LITMUS_PASSWORD = os.Getenv("LITMUS_PASSWORD")
	LITMUS_PROJECT_ID = os.Getenv("LITMUS_PROJECT_ID")

	AGENT_NAME = os.Getenv("AGENT_NAME")
	AGENT_DESCRIPTION = os.Getenv("AGENT_DESCRIPTION")

	NAMESPACE = os.Getenv("NAMESPACE")
	RELEASE_NAME = os.Getenv("RELEASE_NAME")
	APP_VERSION = os.Getenv("APP_VERSION")
	CONFIG_MAP_NAME = os.Getenv("CONFIG_MAP_NAME")
	SERVICE_ACCOUNT_NAME = os.Getenv("SERVICE_ACCOUNT_NAME")

	CLUSTER_ID = os.Getenv("CLUSTER_ID")

}

func GetAgentList(c types.Credentials, project_id string) agentResponse {
	var agentResp agentResponse

	query := `{"query":"query{\n  getCluster(project_id: \"` + LITMUS_PROJECT_ID + `\" ){\n  cluster_id cluster_name access_key is_cluster_confirmed \n  }\n}"}`
	params := apis.SendRequestParams{Endpoint: LITMUS_BACKEND_URL + "/query", Token: c.Token}
	resp, err := apis.SendRequest(params, []byte(query))
	if err != nil {
		utils.Red.Println("Error in getting agent list: ", err)
		os.Exit(1)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		utils.Red.Println("Error in getting agent list: ", err)
		os.Exit(1)
	}

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(bodyBytes, &agentResp)
		if err != nil {
			utils.Red.Println("Error in getting agent list: ", err)
			os.Exit(1)
		}
		return agentResp
	}
	return agentResponse{}
}

func GetAgentWithID(credentials types.Credentials, cluster_id string) agentData {
	agents := GetAgentList(credentials, LITMUS_PROJECT_ID)
	for _, agent := range agents.Data.GetCluster {
		if agent.ClusterID == cluster_id {
			return agent
		}
	}
	return agentData{}

}

func GetAgentWithName(credentials types.Credentials, AgentName string) agentData {
	agents := GetAgentList(credentials, LITMUS_PROJECT_ID)
	for _, agent := range agents.Data.GetCluster {
		if agent.ClusterName == AgentName {
			return agent
		}
	}
	return agentData{}
}

func createConfigMap(configmapName string, configMapData map[string]string) {
	var err error
	config, err := rest.InClusterConfig()
	if err != nil {
                utils.Red.Println("\n❌ Cannot create config from incluster: " + err.Error() + "\n")
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
                utils.Red.Println("\n❌ Cannot create clientset: " + err.Error() + "\n")
	}

	configMap := corev1r.ConfigMap{
		TypeMeta: metav1r.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1r.ObjectMeta{
			Name:      configmapName,
			Namespace: NAMESPACE,
		},
		Data: configMapData,
	}

	var cm *corev1r.ConfigMap
	if _, err := clientset.CoreV1().ConfigMaps(NAMESPACE).Get(context.TODO(), configmapName, metav1r.GetOptions{}); errors.IsNotFound(err) {
		cm, err = clientset.CoreV1().ConfigMaps(NAMESPACE).Create(context.TODO(), &configMap, metav1r.CreateOptions{})
		if err != nil {
			utils.Red.Println("\n❌ Cannot create configmap " + configmapName + " : " + err.Error() + "\n")
			os.Exit(1)
		}

	} else {
		cm, err = clientset.CoreV1().ConfigMaps(NAMESPACE).Update(context.TODO(), &configMap, metav1r.UpdateOptions{})
		if err != nil {
			utils.Red.Println("\n❌ Cannot update configmap " + configmapName + " : " + err.Error() + "\n")
			os.Exit(1)
		}

	}
	_ = cm
}

func createAgent(credentials types.Credentials) {
	var newAgent types.Agent
	newAgent.AgentName = AGENT_NAME
	newAgent.Namespace = NAMESPACE
	newAgent.Description = AGENT_DESCRIPTION
	newAgent.ProjectId = LITMUS_PROJECT_ID
	// Mode
	// 1. cluster
	// 2. namespace
	newAgent.Mode = "cluster"

	// -- OPTIONNAL -- //
	newAgent.ClusterType = "external"
	newAgent.NodeSelector = ""
	// PlatformName
	// 1. AWS
	// 2. GKE
	// 3. Openshift
	// 4. Rancher
	// 5. Others
	newAgent.PlatformName = "Others"
	newAgent.ServiceAccount = SERVICE_ACCOUNT_NAME
	newAgent.SAExists = true
	newAgent.NsExists = true

	t := time.Now()
	configMapData := make(map[string]string)
	configMapData["SERVER_ADDR"] = LITMUS_BACKEND_URL + "/query"
	configMapData["VERSION"] = APP_VERSION
	configMapData["IS_CLUSTER_CONFIRMED"] = "false"
	configMapData["START_TIME"] = t.Format("20060102150405")
	test := `["app.kubernetes.io/instance=` + RELEASE_NAME + `"]`
	configMapData["COMPONENTS"] = "DEPLOYMENTS: " + test
	configMapData["AGENT_SCOPE"] = newAgent.Mode

	var clusterID string
	agentExist := GetAgentWithName(credentials, newAgent.AgentName)
	if (agentExist == agentData{}) {
		utils.White_B.Println("\n🚀 Registering new agent !! 🎉")
		agent, err := apis.ConnectAgent(newAgent, credentials)
		if err != nil {
			utils.Red.Println("\n❌ Agent connection failed: " + err.Error() + "\n")
			os.Exit(1)
		}

		// Print error message in case Data field is null in response
		if (agent.Data == apis.AgentConnect{}) {
			utils.PrintInJsonFormat(agent)
			utils.Red.Println("\n❌ Agent connection failed, null response")
			os.Exit(1)
		}
		clusterID = agent.Data.UserAgentReg.ClusterID
		configMapData["CLUSTER_ID"] = clusterID
		reqCluster := GetAgentWithID(credentials, agent.Data.UserAgentReg.ClusterID)
		if (reqCluster == agentData{}) {
			utils.Red.Println("\n❌ Agent Registered failed: " + err.Error() + "\n")
			os.Exit(1)
		}

		// Checking if cluster with given clusterID and accesskey is present
		configMapData["ACCESS_KEY"] = reqCluster.AccessKey

		utils.White_B.Println("\n🚀 Agent Registered Successful!! 🎉")
	} else {
		clusterID = agentExist.ClusterID

		utils.White_B.Println("\n🚀 Agent Already Registered!! 🎉")

		configMapData["CLUSTER_ID"] = clusterID
		configMapData["ACCESS_KEY"] = agentExist.AccessKey
                configMapData["IS_CLUSTER_CONFIRMED"] = strconv.FormatBool(agentExist.IsClusterConfirmed)

	}

	createConfigMap(CONFIG_MAP_NAME, configMapData)

	configMapWorkflowController := make(map[string]string)
	configMapWorkflowController["config"] = `    containerRuntimeExecutor: k8sapi
    executor:
      imagePullPolicy: IfNotPresent
    instanceID: ` + clusterID
	createConfigMap("workflow-controller-configmap", configMapWorkflowController)

	utils.White_B.Println("\n🚀 Agent Configured Successful!! 🎉")
	utils.White_B.Println("\n🚀 Starting... 🎉")
}

func deleteAgent(credentials types.Credentials) {
        utils.White_B.Println("\n🚀 Delete cluster!! 🎉")

	//query := `{"query":"mutation {\n  deleteClusterReg(clusterInput: \n    { \n    cluster_id: \"` + CLUSTER_ID + `\",\n  }){ cluster_id\n }\n}"}`
	query := `{"operationName":"deleteCluster","variables":{"cluster_id":"` + CLUSTER_ID + `"},"query":"mutation deleteCluster($cluster_id: String\u0021) {\\n  deleteClusterReg(cluster_id: $cluster_id)\\n}\\n"}`
        params := apis.SendRequestParams{Endpoint: LITMUS_BACKEND_URL + "/query", Token: credentials.Token}
	resp, err := apis.SendRequest(params, []byte(query))
	if err != nil {
		utils.Red.Println("Error in getting agent list: ", err)
		os.Exit(1)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		utils.Red.Println("Error in getting agent list: ", err)
		os.Exit(1)
	}
	_ = bodyBytes
        utils.White_B.Println("\n🚀 Agent deleted Successful!! 🎉")

}

func main() {
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
	credentials.Token = resp.AccessToken

	if ACTION == "create" {
                utils.White_B.Println("\n🚀 Start Pre install hook ... 🎉")
		createAgent(credentials)
	} else if ACTION == "delete" {
                utils.White_B.Println("\n🚀 Start Pre delete hook ... 🎉")
		deleteAgent(credentials)
	} else {
		utils.Red.Println("\n❌ Please choose an action, delete or create")
	}
}

// agent.AgentName
// agent.Description
// agent.PlatformName
// agent.ProjectId
// agent.ClusterType
// agent.Mode
// agent.Namespace
// agent.ServiceAccount
// agent.NsExists
//agent.SAExists
