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
	ACTION              string
	LITMUS_FRONTEND_URL string
	LITMUS_USERNAME     string
	LITMUS_PASSWORD     string
)

func init() {
	flag.StringVar(&ACTION, "action", "", "create|delete litmus agent")
	flag.Parse()

	// For all litmus-helm-agent to ChaosCenter communications, This will apply to all requests.
	if os.Getenv("SKIP_SSL") == "true" {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	LITMUS_FRONTEND_URL = os.Getenv("LITMUS_FRONTEND_URL")
	LITMUS_USERNAME = os.Getenv("LITMUS_USERNAME")
	LITMUS_PASSWORD = os.Getenv("LITMUS_PASSWORD")
}

func main() {

	credentials := litmus.Login(LITMUS_FRONTEND_URL, LITMUS_USERNAME, LITMUS_PASSWORD)

	if ACTION == "create" {
		fmt.Println("\n🚀 Start Pre install hook ... 🎉")
		litmus.CreateInfra(credentials)
	} else {
		fmt.Println("\n❌ Please provide a valid action")
	}
}
