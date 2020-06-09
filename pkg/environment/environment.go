package environment

import (
	"os"
	"strconv"

	types "github.com/litmuschaos/test-tools/pkg/types"
	clientTypes "k8s.io/apimachinery/pkg/types"
)

//GetENV fetches all the env variables from the runner pod
func GetENV(experimentDetails *types.ExperimentDetails) {
	experimentDetails.ExperimentName = "pod-delete"
	experimentDetails.ChaosNamespace = os.Getenv("CHAOS_NAMESPACE")
	experimentDetails.EngineName = os.Getenv("CHAOS_ENGINE")
	experimentDetails.ChaosDuration, _ = strconv.Atoi(os.Getenv("TOTAL_CHAOS_DURATION"))
	experimentDetails.Iterations, _ = strconv.Atoi(os.Getenv("ITERATIONS"))
	experimentDetails.ChaosInterval, _ = strconv.Atoi(os.Getenv("CHAOS_INTERVAL"))
	experimentDetails.AppNS = os.Getenv("APP_NS")
	experimentDetails.AppLabel = os.Getenv("APP_LABEL")
	experimentDetails.KillCount, _ = strconv.Atoi(os.Getenv("KILL_COUNT"))
	experimentDetails.ChaosUID = clientTypes.UID(os.Getenv("CHAOS_UID"))
	experimentDetails.ChaosPodName = os.Getenv("POD_NAME")
	experimentDetails.Force, _ = strconv.ParseBool(os.Getenv("FORCE"))
}

//SetEventAttributes initialise all the chaos result ENV
func SetEventAttributes(eventsDetails *types.EventDetails, Reason string, Message string) {

	eventsDetails.Reason = Reason
	eventsDetails.Message = Message
}
