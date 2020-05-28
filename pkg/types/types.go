package types

const (
	// ChaosInject ..
	ChaosInject string = "ChaosInject"
)

// ExperimentDetails is for collecting all the experiment-related details
type ExperimentDetails struct {
	ExperimentName       string
	AppNS                string
	ApplicationContainer string
	ChaosDuration        int
	ChaosInterval        int
	ApplicationPod       string
	Delay                int
	Retry                int
	Iterations           int
}
