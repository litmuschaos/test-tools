package types

import clientTypes "k8s.io/apimachinery/pkg/types"

const (
	// ChaosInject ..
	ChaosInject string = "ChaosInject"
)

// ExperimentDetails is for collecting all the experiment-related details
type ExperimentDetails struct {
	ExperimentName       string
	AppNS                string
	AppLabel             string
	ApplicationContainer string
	ChaosDuration        int
	ChaosInterval        int
	ApplicationPod       string
	Delay                int
	Retry                int
	Iterations           int
	ChaosUID             clientTypes.UID
	ChaosPodName         string
	ChaosNamespace       string
	EngineName           string
	KillCount            int
	Force                bool
}

// EventDetails is for collecting all the events-related details
type EventDetails struct {
	Message string
	Reason  string
}

// ProbeDetails is for collecting all the probe details
type ProbeDetails struct {
	Name                   string
	Type                   string
	Status                 map[string]string
	IsProbeFailedWithError error
	RunID                  string
}

// ResultDetails is for collecting all the chaos-result-related details
type ResultDetails struct {
	Name             string
	Verdict          string
	FailStep         string
	Phase            string
	ResultUID        clientTypes.UID
	ProbeDetails     []ProbeDetails
	PassedProbeCount int
	ProbeArtifacts   map[string]ProbeArtifact
}

// ProbeArtifact contains the probe artifacts
type ProbeArtifact struct {
	ProbeArtifacts RegisterDetails
}

// RegisterDetails contains the output of the corresponding probe
type RegisterDetails struct {
	Register string
}
