package environment

import (
	types "github.com/litmuschaos/test-tools/pkg/types"
)

//GetENV fetches all the env variables from the runner pod
func GetENV(experimentDetails *types.ExperimentDetails) {
	experimentDetails.ExperimentName = "container-kill"
	experimentDetails.AppNS = "test"
	// experimentDetails.ApplicationContainer = os.Getenv("APP_CONTAINER")
	// experimentDetails.ApplicationPod = os.Getenv("APP_POD")
	// experimentDetails.ChaosDuration, _ = strconv.Atoi(os.Getenv("TOTAL_CHAOS_DURATION"))
	// experimentDetails.ChaosInterval, _ = strconv.Atoi(os.Getenv("CHAOS_INTERVAL"))
	experimentDetails.Retry = 90
	experimentDetails.Delay = 2

	experimentDetails.ApplicationContainer = "nginx"
	experimentDetails.ApplicationPod = "nginx-bfb66d6c9-fxtgl	"
	experimentDetails.ChaosDuration = 30
	experimentDetails.ChaosInterval = 10
}
