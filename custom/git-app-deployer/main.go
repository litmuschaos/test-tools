package main

import (
	"bytes"
	"flag"
	"os/exec"
	"strconv"
	"time"

	"github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	"github.com/litmuschaos/litmus-go/pkg/clients"
	"github.com/litmuschaos/litmus-go/pkg/probe"
	types "github.com/litmuschaos/litmus-go/pkg/types"
	"github.com/litmuschaos/test-tools/pkg/log"
	"github.com/openebs/maya/pkg/util/retry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	namespace, filePath, timeout := GetData()

	config, err := getKubeConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	clients := clients.ClientSets{}
	chaosDetails := types.ChaosDetails{}
	resultDetails := types.ResultDetails{}
	experimentLabel := map[string]string{}
	experimentLabel["name"] = resultDetails.Name

	log.Info("[Status]: Starting App Deployer...")
	log.Infof("[Status]: FilePath for App Deployer is %v", filePath)

	if err := CreateNamespace(clientset, namespace); err != nil {
		if k8serrors.IsAlreadyExists(err) {
			chaosResult, err := clients.LitmusClient.ChaosResults(chaosDetails.ChaosNamespace).Get(resultDetails.Name, metav1.GetOptions{})
			if err != nil {
				log.Errorf("Unable to find the chaosresult, err: %v", err)
			}

			// updating the chaosresult with new values
			err = PatchChaosResult(chaosResult, clients, &chaosDetails, &resultDetails, experimentLabel)
			if err != nil {
				log.Errorf("err: %v", err)
			}
		}
		log.Infof("[Status]: %v namespace already exist!", namespace)
	}
	if err := CreateSockShop("/var/run/"+filePath, namespace); err != nil {
		log.Errorf("Failed to install sock-shop, err: %v", err)
		return
	}
	log.Info("[Status]: Sock Shop applications has been successfully created!")

	if err := CheckApplicationStatus(namespace, "app=sock-shop", timeout, 2, clientset); err != nil {
		log.Errorf("err: %v", err)
		return
	}

}

//PatchChaosResult Update the chaos result
func PatchChaosResult(result *v1alpha1.ChaosResult, clients clients.ClientSets, chaosDetails *types.ChaosDetails, resultDetails *types.ResultDetails, chaosResultLabel map[string]string) error {

	result.Status.ExperimentStatus.Phase = resultDetails.Phase
	result.Status.ExperimentStatus.Verdict = resultDetails.Verdict
	result.Spec.InstanceID = chaosDetails.InstanceID
	result.Status.ExperimentStatus.FailStep = resultDetails.FailStep
	// for existing chaos result resource it will patch the label
	result.ObjectMeta.Labels = chaosResultLabel
	result.Status.ProbeStatus = GetProbeStatus(resultDetails)
	if resultDetails.Phase == "Completed" {
		if resultDetails.Verdict == "Pass" && len(resultDetails.ProbeDetails) != 0 {
			result.Status.ExperimentStatus.ProbeSuccessPercentage = "100"

		} else if (resultDetails.Verdict == "Fail" || resultDetails.Verdict == "Stopped") && len(resultDetails.ProbeDetails) != 0 {
			probe.SetProbeVerdictAfterFailure(resultDetails)
			result.Status.ExperimentStatus.ProbeSuccessPercentage = strconv.Itoa((resultDetails.PassedProbeCount * 100) / len(resultDetails.ProbeDetails))
		}

	} else if len(resultDetails.ProbeDetails) != 0 {
		result.Status.ExperimentStatus.ProbeSuccessPercentage = "Awaited"
	}

	// It will update the existing chaos-result CR with new values
	// it will retries until it will able to update successfully or met the timeout(3 mins)
	err := retry.
		Times(90).
		Wait(2 * time.Second).
		Try(func(attempt uint) error {
			_, err := clients.LitmusClient.ChaosResults(result.Namespace).Update(result)
			if err != nil {
				return errors.Errorf("Unable to update the chaosresult, err: %v", err)
			}
			return nil
		})

	return err
}

// GetKubeConfig function derive the kubeconfig
func getKubeConfig() (*rest.Config, error) {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
	// It uses in-cluster config if kubeconfig path is not specified
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	return config, err
}

//GetData derive the sock-shop filePath and timeout
//it derive the filePath based on sock-shop scenario(week vs resilient)
func GetData() (string, string, int) {
	namespace := flag.String("namespace", "", "namespace for the application")
	typeName := flag.String("typeName", "", "type of the application")
	timeout := flag.Int("timeout", 300, "timeout for application status")
	flag.Parse()

	if *namespace == "loadtest" {
		return *namespace, "load-test.yaml", *timeout
	}
	if *typeName == "" || *typeName == "weak" {
		return *namespace, "weak-sock-shop.yaml", *timeout
	}
	return *namespace, *typeName + "-sock-shop.yaml", *timeout
}

// CreateNamespace creates a sock-shop namespace
func CreateNamespace(clientset *kubernetes.Clientset, namespaceName string) error {
	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
	_, err := clientset.CoreV1().Namespaces().Create(nsSpec)
	return err
}

//CreateSockShop creates sock-shop application
func CreateSockShop(path string, ns string) error {
	command := exec.Command("kubectl", "apply", "-f", path, "-n", ns)
	var out, stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		log.Infof(" %v", stderr.String())
		return err
	}
	return nil
}

// CheckApplicationStatus checks the status of the AUT
func CheckApplicationStatus(appNs, appLabel string, timeout, delay int, clientset *kubernetes.Clientset) error {
	// Checking whether application containers are in ready state
	log.Info("[Status]: Checking whether application containers are in ready state")
	err := CheckContainerStatus(appNs, appLabel, timeout, delay, clientset)
	if err != nil {
		return err
	}
	// Checking whether application pods are in running state
	log.Info("[Status]: Checking whether application pods are in running state")
	err = CheckPodStatus(appNs, appLabel, timeout, delay, clientset)
	if err != nil {
		return err
	}
	return nil
}

// CheckPodStatus checks the running status of the application pod
func CheckPodStatus(appNs, appLabel string, timeout, delay int, clientset *kubernetes.Clientset) error {
	err := retry.
		Times(uint(timeout / delay)).
		Wait(time.Duration(delay) * time.Second).
		Try(func(attempt uint) error {
			podSpec, err := clientset.CoreV1().Pods(appNs).List(metav1.ListOptions{LabelSelector: appLabel})
			if err != nil || len(podSpec.Items) == 0 {
				return errors.Errorf("Unable to find the pods with matching labels, err: %v", err)
			}
			for _, pod := range podSpec.Items {
				if string(pod.Status.Phase) != "Running" {
					return errors.Errorf("Pod is not yet in running state")
				}
				log.InfoWithValues("[Status]: The running status of Pods are as follows", logrus.Fields{
					"Pod": pod.Name, "Status": pod.Status.Phase})
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

// CheckContainerStatus checks the status of the application container
func CheckContainerStatus(appNs, appLabel string, timeout, delay int, clientset *kubernetes.Clientset) error {
	err := retry.
		Times(uint(timeout / delay)).
		Wait(time.Duration(delay) * time.Second).
		Try(func(attempt uint) error {
			podSpec, err := clientset.CoreV1().Pods(appNs).List(metav1.ListOptions{LabelSelector: appLabel})
			if err != nil || len(podSpec.Items) == 0 {
				return errors.Errorf("Unable to find the pods with matching labels, err: %v", err)
			}
			for _, pod := range podSpec.Items {
				for _, container := range pod.Status.ContainerStatuses {
					if container.State.Terminated != nil {
						return errors.Errorf("container is in terminated state")
					}
					if container.Ready != true {
						return errors.Errorf("containers are not yet in ready state")
					}
					log.InfoWithValues("[Status]: The running status of container are as follows", logrus.Fields{
						"container": container.Name, "Pod": pod.Name, "Readiness": container.Ready})
				}
			}
			return nil
		})
	if err != nil {
		return err
	}
	return nil
}

//GetProbeStatus fetch status of all probes
func GetProbeStatus(resultDetails *types.ResultDetails) []v1alpha1.ProbeStatus {

	probeStatus := []v1alpha1.ProbeStatus{}
	for _, probe := range resultDetails.ProbeDetails {
		probes := v1alpha1.ProbeStatus{}
		probes.Name = probe.Name
		probes.Type = probe.Type
		probes.Status = probe.Status
		probeStatus = append(probeStatus, probes)
	}
	return probeStatus
}
