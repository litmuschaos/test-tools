package server

type ChaosType string

const (
	Error            ChaosType = "error"
	RandomResolution ChaosType = "random"
)

type MatchType string

const (
	Exact  MatchType = "exact"
	Substr MatchType = "substring"
)

type InterceptorSettings struct {
	TargetHostNames []string
	MatchType
	ChaosType
}
