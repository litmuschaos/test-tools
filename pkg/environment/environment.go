package environment

import (
	"os"
	"strconv"

	types "github.com/litmuschaos/test-tools/pkg/types"
	clientTypes "k8s.io/apimachinery/pkg/types"
)

//GetENV fetches all the env variables from the runner pod
func GetENV(experimentDetails *types.ExperimentDetails, name string) {
	experimentDetails.ExperimentName = name
	experimentDetails.AppNS = os.Getenv("APP_NS")
	experimentDetails.Retry, _ = strconv.Atoi(Getenv("RETRY", "90"))
	experimentDetails.Delay, _ = strconv.Atoi(Getenv("DELAY", "2"))
	experimentDetails.ApplicationContainer = os.Getenv("APP_CONTAINER")
	experimentDetails.ApplicationPod = os.Getenv("APP_POD")
	experimentDetails.ChaosDuration, _ = strconv.Atoi(Getenv("TOTAL_CHAOS_DURATION", "30"))
	experimentDetails.ChaosInterval, _ = strconv.Atoi(Getenv("CHAOS_INTERVAL", "10"))
	experimentDetails.Iterations, _ = strconv.Atoi(Getenv("ITERATIONS", "3"))
	experimentDetails.ChaosNamespace = Getenv("CHAOS_NAMESPACE", "litmus")
	experimentDetails.EngineName = os.Getenv("CHAOS_ENGINE")
	experimentDetails.AppLabel = os.Getenv("APP_LABEL")
	experimentDetails.KillCount, _ = strconv.Atoi(Getenv("KILL_COUNT", "1"))
	experimentDetails.ChaosUID = clientTypes.UID(os.Getenv("CHAOS_UID"))
	experimentDetails.ChaosPodName = os.Getenv("POD_NAME")
	experimentDetails.Force, _ = strconv.ParseBool(Getenv("FORCE", "false"))
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
