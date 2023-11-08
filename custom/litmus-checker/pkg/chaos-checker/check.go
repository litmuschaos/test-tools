package chaos_checker

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdsoumya/resourceChecker/pkg/k8s"
)

func CheckChaos(kubeconfig *string, res k8s.ResourceDef) {
	dc, dyn, err := k8s.GetDynamicClient(kubeconfig)
	if err != nil {
		log.Fatal("ERROR : ", err)
	}

	checkerInterval := int64(60)

	// poll interval in seconds
	intervalEnv := os.Getenv("CHECK_INTERVAL")
	if intervalEnv != "" {
		checkerInterval, err = strconv.ParseInt(intervalEnv, 10, 64)
		if err != nil {
			log.Fatal("ERROR failed to parse checker interval seconds: ", err)
		}
	}

	log.Printf("Starting Chaos Checker in %v seconds", checkerInterval)

	for {
		time.Sleep(time.Second * time.Duration(checkerInterval))
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
