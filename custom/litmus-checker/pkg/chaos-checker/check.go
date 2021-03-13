package chaos_checker

import (
	"log"
	"strings"
	"time"

	"github.com/gdsoumya/resourceChecker/pkg/k8s"
)

func CheckChaos(kubeconfig *string, res k8s.ResourceDef) {
	dc, dyn, err := k8s.GetDynamicClient(kubeconfig)
	if err != nil {
		log.Fatal("ERROR : ", err)
	}
	log.Print("Starting Chaos Checker in 1min")
	for {
		time.Sleep(time.Minute * 1)
		log.Print("Checking if Engine Completed or Stopped")
		data, err := k8s.GetResourceDetails(dc, dyn, res)
		if err != nil {
			log.Fatal("ERROR : ", err)
		}
		status, ok := data.Object["status"].(map[string]interface{})
		if !ok {
			continue
		}
		engStat, ok := status["engineStatus"].(string)
		if ok {
			if strings.ToLower(engStat) == "completed" {
				log.Print("[*] ENGINE COMPLETED")
				return
			} else if strings.ToLower(engStat) == "stopped" {
				log.Print("[!] ERROR : ENGINE STATUS STOPPED")
				return
			}
		}
	}
}
