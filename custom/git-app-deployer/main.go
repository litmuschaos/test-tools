package main

import (
	"bytes"
	"flag"
	"os/exec"
	"time"

	"github.com/litmuschaos/test-tools/pkg/log"
	"github.com/openebs/maya/pkg/util/retry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//APPVars maintaining all parameters during application installation and uninstallation
type APPVars struct {
	namespace string
	filePath  string
	timeout   int
	operation string
}

func main() {

	appVars := GetData()

	config, err := getKubeConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	log.Info("[Status]: Starting App Deployer...")

	//operations for application
	switch appVars.operation {
	case "apply", "create":
		if err := CreateApplication(appVars, 2, clientset); err != nil {
			log.Errorf("err: %v", err)
			return
		}
		log.Info("[Status]: Sock Shop applications has been successfully created!")
	case "delete":
		if err := DeleteApplication(appVars, 2, clientset); err != nil {
			log.Errorf("err: %v", err)
			return
		}
		log.Info("[Status]: Sock Shop applications has been successfully deleted!")
	default:
		log.Infof("operation '%s' not supported in app-deployer", appVars.operation)
		return
	}
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
func GetData() (vars *APPVars) {

	//Initialise the variables
	appVars := APPVars{}

	namespace := flag.String("namespace", "", "namespace for the application")
	filePath := flag.String("typeName", "weak", "type of the application")
	timeout := flag.Int("timeout", 300, "timeout for application status")
	operation := flag.String("operation", "apply", "type of operation for application")
	flag.Parse()

	appVars.namespace = *namespace
	appVars.timeout = *timeout
	appVars.operation = *operation

	//For sock-shop namespace weak/resilient filePath is
	//required and for loadtest namespace load-test filePath
	switch appVars.namespace {
	case "loadtest":
		appVars.filePath = "load-test.yaml"
	case "sock-shop":
		appVars.filePath = *filePath + "-sock-shop.yaml"
	default:
		log.Infof("namespace '%s' not supported in app-deployer", appVars.operation)
		return
	}

	return &appVars
}

// CreateNamespace creates a namespace
func CreateNamespace(clientset *kubernetes.Clientset, namespaceName string) error {
	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespaceName}}
	_, err := clientset.CoreV1().Namespaces().Create(nsSpec)
	return err
}

// DeleteNamespace delete a namespace
func DeleteNamespace(clientset *kubernetes.Clientset, namespaceName string) error {
	err := clientset.CoreV1().Namespaces().Delete(namespaceName, &metav1.DeleteOptions{})
	return err
}

//CreateSockShop install the application
func CreateSockShop(path string, ns string, operation string) error {
	command := exec.Command("kubectl", operation, "-f", path, "-n", ns)
	var out, stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		log.Infof(" %v", stderr.String())
		return err
	}
	return nil
}

//DeleteSockShop uninstall the application
func DeleteSockShop(path string, ns string, operation string) error {
	command := exec.Command("kubectl", operation, "-f", path, "-n", ns)
	var out, stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		log.Infof(" %v", stderr.String())
		return err
	}
	return nil
}

//CreateApplication install the application and add  all corresponding resources
func CreateApplication(appVars *APPVars, delay int, clientset *kubernetes.Clientset) error {

	log.Infof("[Status]: FilePath for App Deployer is %v", appVars.filePath)

	if err := CreateNamespace(clientset, appVars.namespace); err != nil {
		log.Info("[Status]: Namespace already exist!")
	}

	if err := CreateSockShop("/var/run/"+appVars.filePath, appVars.namespace, appVars.operation); err != nil {
		log.Errorf("Failed to install sock-shop, err: %v", err)
		return err
	}

	if err := CheckApplicationStatus(appVars.namespace, "app=sock-shop", appVars.timeout, 2, appVars.operation, clientset); err != nil {
		log.Errorf("err: %v", err)
		return err
	}

	return nil
}

//DeleteApplication deletes the application and remove all corresponding resources
func DeleteApplication(appVars *APPVars, delay int, clientset *kubernetes.Clientset) error {

	log.Infof("[Status]: FilePath for App Deployer is %v", appVars.filePath)

	if err := DeleteSockShop("/var/run/"+appVars.filePath, appVars.namespace, appVars.operation); err != nil {
		return err
	}

	if err := CheckPodStatusForRevert(appVars.namespace, appVars.timeout, 2, clientset); err != nil {
		return err
	}

	if err := DeleteNamespace(clientset, appVars.namespace); err != nil {
		log.Info("[Status]: Namespace not found!")
	}

	return nil
}

// CheckApplicationStatus checks the status of the AUT
func CheckApplicationStatus(appNs, appLabel string, timeout, delay int, operation string, clientset *kubernetes.Clientset) error {
	//  Checking application containers state
	log.Info("[Status]: Checking application containers state")
	err := CheckContainerStatus(appNs, appLabel, timeout, delay, clientset)
	if err != nil {
		return err
	}
	// Checking application pods state
	log.Info("[Status]: Checking application pods state")
	err = CheckPodStatus(appNs, appLabel, timeout, delay, clientset)
	if err != nil {
		return err
	}
	return nil
}

//CheckPodStatusForRevert wait for the application to terminate all pods
func CheckPodStatusForRevert(appNs string, timeout, delay int, clientset *kubernetes.Clientset) error {
	err := retry.
		Times(uint(timeout / delay)).
		Wait(time.Duration(delay) * time.Second).
		Try(func(attempt uint) error {
			podSpec, err := clientset.CoreV1().Pods(appNs).List(metav1.ListOptions{})
			if err != nil {
				return errors.Errorf("Unable to find the pods in namespace, err: %v", err)
			}

			if len(podSpec.Items) != 0 {
				return errors.Errorf("[Status]: Pods are yet to be terminated")
			}
			return nil
		})
	return err
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
