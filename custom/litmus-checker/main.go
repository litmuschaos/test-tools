package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	chaos_checker "github.com/litmuschaos/test-tools/custom/litmus-checker/pkg/chaos-checker"
	"github.com/litmuschaos/test-tools/custom/litmus-checker/pkg/k8s"
	"github.com/litmuschaos/test-tools/custom/litmus-checker/pkg/util"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	file := flag.String("file", "", "absolute path to the chaosengine yaml")
	engineFile := flag.String("saveName", "", "absolute path to the output file")
	flag.Parse()

	if file == nil || *file == "" {
		logrus.Fatal("Error Engine Artefact path not specified")
	}

	data, err := ioutil.ReadFile(*file)
	if err != nil {
		logrus.Fatalf("Error Reading Artefact : %v", err)
	}

	dc, dyn, err := k8s.GetDynamicClient(kubeconfig)
	if err != nil {
		logrus.Fatalf("Error Getting Dynamic Client : %v", err)
	}

	resp, err := k8s.CreateChaosDeployment(data, dc, dyn)
	if err != nil {
		logrus.Fatalf("Error Creating Resource : %v", err)
	}

	engineName := resp.GetName()
	logrus.Infof("ChaosEngine Name : %s", engineName)

	if err = util.WriteToFile(*engineFile, engineName); err != nil {
		logrus.Infof("ERROR: cannot write engine-name  %v", err)
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
		logrus.Info("SIGTERM SIGNAL RECEIVED, Shutting down litmus-checker...")
		os.Exit(0)
	}()

	logrus.Infof("Created Resource Details: %v", resDef)
	chaos_checker.CheckChaos(resDef, dc, dyn)
}
