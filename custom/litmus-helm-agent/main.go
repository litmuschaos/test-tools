package main

import (
	"flag"
	"fmt"
	litmus "litmus-helm-agent/pkg/litmus"
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

	LITMUS_FRONTEND_URL = os.Getenv("LITMUS_FRONTEND_URL")
	LITMUS_USERNAME = os.Getenv("LITMUS_USERNAME")
	LITMUS_PASSWORD = os.Getenv("LITMUS_PASSWORD")
}

func main() {

	credentials := litmus.Login(LITMUS_FRONTEND_URL, LITMUS_USERNAME, LITMUS_PASSWORD)

	if ACTION == "create" {
		fmt.Println("\n🚀 Start Pre install hook ... 🎉")
		litmus.CreateAgent(credentials)

	} else if ACTION == "delete" {
		fmt.Println("\n🚀 Start Pre delete hook ... 🎉")
		litmus.DeleteAgent(credentials)

	} else {
		fmt.Println("\n❌ Please choose an action, delete or create")

	}
}
