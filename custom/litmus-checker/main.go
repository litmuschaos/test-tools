package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	chaos_checker "github.com/gdsoumya/resourceChecker/pkg/chaos-checker"
	"github.com/gdsoumya/resourceChecker/pkg/k8s"
	"github.com/gdsoumya/resourceChecker/pkg/util"
)

func main() {
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	file := flag.String("file", "", "absolute path to the chaosengine yaml")
	engineFile := flag.String("saveName", "", "absolute path to the output file")
	flag.Parse()
	if *file == "" {
		log.Fatal("Error Engine Artefact path not specified")
	}
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal("Error Reading Artefact : ", err)
	}
	resp, err := k8s.CreateChaosDeployment(data, kubeconfig)
	if err != nil {
		log.Fatal("Error Creating Resource : ", err)
	}
	engineName := resp.GetName()
	log.Print("\n\nChaosEngine Name : ", engineName, "\n\n")
	err = util.WriteToFile(*engineFile, engineName)
	if err != nil {
		log.Print("ERROR : cannot write engine-name - ", err)
	}
	gvk := resp.GroupVersionKind()
	resDef := k8s.ResourceDef{
		Name:      engineName,
		Group:     gvk.Group,
		Version:   gvk.Version,
		Kind:      gvk.Kind,
		Namespace: resp.GetNamespace(),
		Selectors: "",
	}

	// Required, While aborting a Chaos Experiment, wait-container (argo-exec) sends SIGTERM signal to other (main) containers for aborting Argo-Workflow Pod
	go func() {
		<-signalChannel
		log.Print("SIGTERM SIGNAL RECEIVED, Shutting down litmus-checker...")
		os.Exit(0)
	}()

	log.Print("Created Resource Details: \n", resDef)
	chaos_checker.CheckChaos(kubeconfig, resDef)
}
