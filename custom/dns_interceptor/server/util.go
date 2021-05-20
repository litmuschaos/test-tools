package server

import (
	"bufio"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// getInterceptorSettings generates the interceptor settings from the env
func getInterceptorSettings() (*InterceptorSettings, error) {
	targets := os.Getenv("TARGET_HOSTNAMES")
	spoofMap := os.Getenv("SPOOF_MAP")
	chaosType := os.Getenv("CHAOS_TYPE")
	matchType := os.Getenv("MATCH_SCHEME")
	interceptorSettings := InterceptorSettings{}

	if targets == "" {
		interceptorSettings.TargetHostNames = nil
	} else {
		err := json.Unmarshal([]byte(targets), &interceptorSettings.TargetHostNames)
		if err != nil {
			return nil, errors.New("failed to parse target hostname list : " + err.Error())
		}
	}

	if spoofMap == "" {
		interceptorSettings.SpoofMap = nil
	} else {
		err := json.Unmarshal([]byte(spoofMap), &interceptorSettings.SpoofMap)
		if err != nil {
			return nil, errors.New("failed to parse target hostname list : " + err.Error())
		}
	}

	log.WithField("spoof_map", interceptorSettings.SpoofMap).Info("Chaos Spoof Map")
	//currently only Error type is supported
	if ChaosType(chaosType) != Error && ChaosType(chaosType) != Spoof {
		return nil, errors.New("wrong chaos type for dns chaos")
	} else {
		interceptorSettings.ChaosType = ChaosType(chaosType)
	}

	log.WithField("chaos-type", interceptorSettings.ChaosType).Info("Chaos type")

	// defaults to Exact
	if matchType == "" {
		interceptorSettings.MatchType = Exact
	} else if MatchType(matchType) != Exact && MatchType(matchType) != Substr {
		return nil, errors.New("wrong chaos type for dns chaos")
	} else {
		interceptorSettings.MatchType = MatchType(matchType)
	}

	log.WithField("match-scheme", interceptorSettings.MatchType).Info("Target match scheme")
	return &interceptorSettings, nil
}

// updateResolvConf updates the resolv.conf file with the require dns nameserver
func updateResolvConf(resolvConfPath string, originalData *string) (string, error) {
	file, err := os.OpenFile(resolvConfPath, os.O_RDWR, 644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	newLines := ""
	originalLines := ""
	if originalData == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "nameserver") {
				newLines += "nameserver 127.0.0.1\n"
			} else {
				newLines += line + "\n"
			}
			originalLines += line + "\n"
		}
	} else {
		newLines = *originalData
	}
	err = file.Truncate(0)
	if err != nil {
		return "", err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	_, err = file.WriteString(newLines)
	if err != nil {
		return "", err
	}
	return originalLines, nil
}
