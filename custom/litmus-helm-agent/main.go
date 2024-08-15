package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	litmus "litmus-helm-agent/pkg/litmus"
	"net/http"
	"os"
)

var (
	ACTION     string
	INFRA_ID   string
	ACCESS_KEY string
)

func init() {
	flag.StringVar(&ACTION, "action", "", "create|delete litmus agent")
	flag.Parse()

	// For all litmus-helm-agent to ChaosCenter communications, This will apply to all requests.
	if os.Getenv("SKIP_SSL") == "true" {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	INFRA_ID = os.Getenv("INFRA_ID")
	ACCESS_KEY = os.Getenv("ACCESS_KEY")
}

func main() {

	if ACTION == "create" {
		fmt.Println("\n🚀 Start Pre install hook ... 🎉")
		litmus.CreateInfra(INFRA_ID, ACCESS_KEY)
	} else {
		fmt.Println("\n❌ Please provide a valid action")
	}
}
