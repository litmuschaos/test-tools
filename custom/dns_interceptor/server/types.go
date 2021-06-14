package server

type ChaosType string

const (
	Error            ChaosType = "error"
	Spoof            ChaosType = "spoof"
	RandomResolution ChaosType = "random"
	DefaultDNSPort             = "53"
)

type MatchType string

const (
	Exact  MatchType = "exact"
	Substr MatchType = "substring"
)

type InterceptorSettings struct {
	TargetHostNames []string
	SpoofMap        map[string]string
	MatchType
	ChaosType
}
