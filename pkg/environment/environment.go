package environment

import (
	"os"
	"strconv"

	types "github.com/litmuschaos/test-tools/pkg/types"
	clientTypes "k8s.io/apimachinery/pkg/types"
)

//GetENV fetches all the env variables from the runner pod
func GetENV(experimentDetails *types.ExperimentDetails) {
	experimentDetails.ExperimentName = "container-kill"
	experimentDetails.AppNS = Getenv("APP_NS", "test")
	experimentDetails.Retry = 90
	experimentDetails.Delay = 2
	experimentDetails.ApplicationContainer = Getenv("APP_CONTAINER", "nginx")
	experimentDetails.ApplicationPod = Getenv("APP_POD", "nginx-bfb66d6c9-x5gks")
	experimentDetails.ChaosDuration, _ = strconv.Atoi(Getenv("TOTAL_CHAOS_DURATION", "30"))
	experimentDetails.ChaosInterval, _ = strconv.Atoi(Getenv("CHAOS_INTERVAL", "10"))
	experimentDetails.Iterations, _ = strconv.Atoi(Getenv("CHAOS_ITERATION", "3"))
	experimentDetails.ChaosNamespace = Getenv("CHAOS_NS", "test")
	experimentDetails.EngineName = Getenv("ENGINE_NAME", "")
	experimentDetails.ChaosUID = clientTypes.UID(Getenv("ENGINE_UID", ""))
	experimentDetails.ChaosPodName = Getenv("CHAOS_POD", "test")
}

//SetEventAttributes initialise all the chaos result ENV
func SetEventAttributes(eventsDetails *types.EventDetails, Reason string, Message string) {

	eventsDetails.Reason = Reason
	eventsDetails.Message = Message
}

// Getenv fetch the env and set the default value, if any
func Getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
