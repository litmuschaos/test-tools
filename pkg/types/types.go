package types

import clientTypes "k8s.io/apimachinery/pkg/types"

const (
	// ChaosInject ..
	ChaosInject string = "ChaosInject"
)

// ExperimentDetails is for collecting all the experiment-related details
type ExperimentDetails struct {
	ExperimentName string
	AppNS          string
	AppLabel       string
	EngineName     string
	ChaosNamespace string
	ChaosUID       clientTypes.UID
	ChaosPodName   string
	ChaosDuration  int
	ChaosInterval  int
	Iterations     int
	KillCount      int
	Force          bool
}

// EventDetails is for collecting all the events-related details
type EventDetails struct {
	Message string
	Reason  string
}
