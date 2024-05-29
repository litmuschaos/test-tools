package chaos_checker

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"strings"
	"time"

	"github.com/litmuschaos/test-tools/custom/litmus-checker/pkg/k8s"
)

func CheckChaos(res k8s.ResourceDef, dc discovery.DiscoveryInterface, dyn dynamic.Interface) {
	logrus.Info("Starting Chaos Checker in 1min")

	for {
		time.Sleep(time.Minute * 1)
		data, err := k8s.GetResourceDetails(dc, dyn, res)
		if err != nil {
			logrus.Fatalf("Failed to get resource details: %v", err)
		}
		status, ok := data.Object["status"].(map[string]interface{})
		if !ok {
			continue
		}
		engStat, ok := status["engineStatus"].(string)
		if ok {
			logrus.Infof("Engine Status :%s", engStat)
			if strings.ToLower(engStat) == "completed" {
				logrus.Info("[*] ENGINE COMPLETED")
				return
			} else if strings.ToLower(engStat) == "stopped" {
				logrus.Info("[!] ERROR : ENGINE STATUS STOPPED")
				return
			}
		}
	}
}
