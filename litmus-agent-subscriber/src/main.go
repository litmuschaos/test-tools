package main

import (
  "fmt"
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
	"errors"

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
  AGENT_MODE        string
  AGENT_NODE_SELECTOR string
  
	NAMESPACE            string
	APP_VERSION          string
	SERVICE_ACCOUNT_NAME string
	AGENT_CONFIGMAP_NAME string
	RELEASE_NAME         string
	ACTION               string

	CLUSTER_ID	     string
	CLUSTER_TYPE     string
	PLATFORM_NAME    string
	AGENT_SA_EXISTS        bool
	AGENT_NS_EXISTS        bool
	
	WORKFLOW_CONTROLER_CONFIGMAP_NAME string
	CONTAINER_RUNTIME_EXECUTOR string
	
	
)

type AgentConnectionData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data AgentConnect `json:"data"`
}

type Errors struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type AgentConnect struct {
	UserAgentReg UserAgentReg `json:"userClusterReg"`
}

type UserAgentReg struct {
	ClusterID   string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Token       string `json:"token"`
}
	
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
	AGENT_MODE = os.Getenv("AGENT_MODE")
  AGENT_CONFIGMAP_NAME = os.Getenv("AGENT_CONFIGMAP_NAME")
  AGENT_NODE_SELECTOR = os.Getenv("AGENT_NODE_SELECTOR")
  CLUSTER_TYPE = os.Getenv("CLUSTER_TYPE")
  
	NAMESPACE = os.Getenv("NAMESPACE")
	RELEASE_NAME = os.Getenv("RELEASE_NAME")
	APP_VERSION = os.Getenv("APP_VERSION")

	SERVICE_ACCOUNT_NAME = os.Getenv("SERVICE_ACCOUNT_NAME")

  PLATFORM_NAME = os.Getenv("PLATFORM_NAME")
  AGENT_SA_EXISTS, _ = strconv.ParseBool(os.Getenv("SA_EXISTS"))
  AGENT_NS_EXISTS, _ = strconv.ParseBool(os.Getenv("NS_EXISTS"))
  
  WORKFLOW_CONTROLER_CONFIGMAP_NAME = os.Getenv("WORKFLOW_CONTROLER_CONFIGMAP_NAME")
  CONTAINER_RUNTIME_EXECUTOR = os.Getenv("CONTAINER_RUNTIME_EXECUTOR")

	
	CLUSTER_ID = os.Getenv("CLUSTER_ID")

}

func GetAgentList(c types.Credentials, project_id string) agentResponse {
	var agentResp agentResponse

	query := `{"query":"query{\n  getCluster(project_id: \"` + LITMUS_PROJECT_ID + `\" ){\n  cluster_id cluster_name access_key is_cluster_confirmed \n  }\n}"}`
	params := apis.SendRequestParams{Endpoint: LITMUS_BACKEND_URL + "/query", Token: c.Token}
	resp, err := apis.SendRequest(params, []byte(query), "POST")
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

// ConnectAgent connects the agent with the given details
func ConnectAgent(agent types.Agent, cred types.Credentials) (apis.AgentConnectionData, error) {
	query := `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  \tplatform_name: \"` + agent.PlatformName + `\",\n    project_id: \"` + agent.ProjectId + `\",\n    cluster_type: \"` + agent.ClusterType + `\",\n  agent_scope: \"` + agent.Mode + `\",\n    agent_namespace: \"` + agent.Namespace + `\",\n    serviceaccount: \"` + agent.ServiceAccount + `\",\n    skip_ssl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    agent_ns_exists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n  }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`

	if agent.NodeSelector != "" {
		query = `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  node_selector: \"` + agent.NodeSelector + `\",\n  \tplatform_name: \"` + agent.PlatformName + `\",\n    project_id: \"` + agent.ProjectId + `\",\n    cluster_type: \"` + agent.ClusterType + `\",\n  agent_scope: \"` + agent.Mode + `\",\n    agent_namespace: \"` + agent.Namespace + `\",\n    skip_ssl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    serviceaccount: \"` + agent.ServiceAccount + `\",\n    agent_ns_exists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n  }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`
	}

	if agent.Tolerations != "" {
		query = `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  \tplatform_name: \"` + agent.PlatformName + `\",\n    project_id: \"` + agent.ProjectId + `\",\n    cluster_type: \"` + agent.ClusterType + `\",\n  agent_scope: \"` + agent.Mode + `\",\n    agent_namespace: \"` + agent.Namespace + `\",\n    serviceaccount: \"` + agent.ServiceAccount + `\",\n    skip_ssl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    agent_ns_exists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n tolerations: ` + agent.Tolerations + ` }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`
	}

	if agent.NodeSelector != "" && agent.Tolerations != "" {
		query = `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  node_selector: \"` + agent.NodeSelector + `\",\n  \tplatform_name: \"` + agent.PlatformName + `\",\n    project_id: \"` + agent.ProjectId + `\",\n    cluster_type: \"` + agent.ClusterType + `\",\n  agent_scope: \"` + agent.Mode + `\",\n    agent_namespace: \"` + agent.Namespace + `\",\n    skip_ssl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    serviceaccount: \"` + agent.ServiceAccount + `\",\n    agent_ns_exists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n tolerations: ` + agent.Tolerations + ` }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`
	}

	resp, err := apis.SendRequest(apis.SendRequestParams{Endpoint: LITMUS_BACKEND_URL + "/query", Token: cred.Token}, []byte(query), "POST")
	if err != nil {
		return apis.AgentConnectionData{}, errors.New("Error in registering agent: " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return apis.AgentConnectionData{}, errors.New("Error in registering agent: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var connectAgent apis.AgentConnectionData
		err = json.Unmarshal(bodyBytes, &connectAgent)
		if err != nil {
			return apis.AgentConnectionData{}, errors.New("Error in registering agent: " + err.Error())
		}

		if len(connectAgent.Errors) > 0 {
			return apis.AgentConnectionData{}, errors.New(connectAgent.Errors[0].Message)
		}
		return connectAgent, nil
	} else {
		return apis.AgentConnectionData{}, err
	}
}

func createAgent(credentials types.Credentials) {
	var newAgent types.Agent
	newAgent.AgentName = AGENT_NAME
	newAgent.Namespace = NAMESPACE
	newAgent.Description = AGENT_DESCRIPTION
	newAgent.ProjectId = LITMUS_PROJECT_ID
	newAgent.Mode = AGENT_MODE
	newAgent.SkipSSL = true

	// -- OPTIONNAL -- //
	newAgent.ClusterType = CLUSTER_TYPE
	newAgent.NodeSelector = AGENT_NODE_SELECTOR
	newAgent.PlatformName = PLATFORM_NAME
	newAgent.ServiceAccount = SERVICE_ACCOUNT_NAME
	newAgent.SAExists = AGENT_SA_EXISTS
	newAgent.NsExists = AGENT_NS_EXISTS

	configMapData := make(map[string]string)
	configMapData["SERVER_ADDR"] = LITMUS_BACKEND_URL + "/query"
	configMapData["VERSION"] = APP_VERSION
	configMapData["IS_CLUSTER_CONFIRMED"] = "false"
	configMapData["START_TIME"] = strconv.FormatInt(time.Now().Unix(), 10)
	selector := `["app.kubernetes.io/instance=` + RELEASE_NAME + `"]`
	configMapData["COMPONENTS"] = "DEPLOYMENTS: " + selector
	configMapData["AGENT_SCOPE"] = newAgent.Mode

	var clusterID string
	agentExist := GetAgentWithName(credentials, newAgent.AgentName)
	
	if (agentExist == agentData{}) {
	
		utils.White_B.Println("\n🚀 Registering new agent !! 🎉")
		
		fmt.Printf("%+v\n", newAgent)
		agent, err := ConnectAgent(newAgent, credentials)
		if err != nil {
			utils.Red.Println("\n❌ Agent connection failed: " + err.Error() + "\n")
			os.Exit(1)
		}

		// Print error message in case Data field is null in response
		if (agent.Data == apis.AgentConnect{}) {
			utils.PrintInJsonFormat(agent)
			utils.Red.Println("\n❌ Agent connection failed: unknown error\n")
			os.Exit(1)
		}
		clusterID = agent.Data.UserAgentReg.ClusterID
		configMapData["CLUSTER_ID"] = clusterID
		
		reqCluster := GetAgentWithID(credentials, agent.Data.UserAgentReg.ClusterID)
		if (reqCluster == agentData{}) {
			utils.Red.Println("\n❌ Agent Registered failed\n")
			os.Exit(1)
		}

		// Checking if cluster with given clusterID and accesskey is present
		configMapData["ACCESS_KEY"] = reqCluster.AccessKey

    utils.White_B.Println("\n🚀 Agent Registered Successful!! 🎉")
    createConfigMap(AGENT_CONFIGMAP_NAME, configMapData)
    
    configMapWorkflowController := make(map[string]string)
    configMapWorkflowController["config"] = `    containerRuntimeExecutor: ` + CONTAINER_RUNTIME_EXECUTOR + `
    executor:
      imagePullPolicy: IfNotPresent
    instanceID: ` + clusterID
    
    createConfigMap(WORKFLOW_CONTROLER_CONFIGMAP_NAME, configMapWorkflowController)

    utils.White_B.Println("\n🚀 Agent Configured Successful!! 🎉")
    utils.White_B.Println("\n🚀 Starting... 🎉")
	
	}
}


func deleteAgent(credentials types.Credentials) {
        utils.White_B.Println("\n🚀 Delete cluster!! 🎉")

	//query := `{"query":"mutation {\n  deleteClusterReg(clusterInput: \n    { \n    cluster_id: \"` + CLUSTER_ID + `\",\n  }){ cluster_id\n }\n}"}`
	query := `{"operationName":"deleteCluster","variables":{"cluster_id":"` + CLUSTER_ID + `"},"query":"mutation deleteCluster($cluster_id: String\u0021) {\\n  deleteClusterReg(cluster_id: $cluster_id)\\n}\\n"}`
        params := apis.SendRequestParams{Endpoint: LITMUS_BACKEND_URL + "/query", Token: credentials.Token}
	resp, err := apis.SendRequest(params, []byte(query), "POST")
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
