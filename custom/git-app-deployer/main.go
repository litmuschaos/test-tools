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

	log.Info("[Status]: Starting App Deployer...")
	log.Infof("[Status]: FilePath for App Deployer is %v", filePath)

	InstallationApplication(filePath, namespace, timeout, clientset)
	log.Info("[Status]: Sock Shop applications has been successfully created!")

	InstallationJaeger(timeout, clientset)
	log.Info("[Status]: Jaeger has been successfully created!")

	log.Info("[Status]: Application Is Running !")
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

//InstallationApplication is creating and checking status of sock-shop application
func InstallationApplication(path string, namespace string, timeout int, clientset *kubernetes.Clientset) {
	if err := CreateNamespace(clientset, namespace); err != nil {
		log.Info("[Status]: Namespace already exist!")
	}

	if err := CreateSockShop("/var/run/"+path, namespace); err != nil {
		log.Errorf("Failed to install sock-shop, err: %v", err)
		return
	}
	if err := CheckApplicationStatus(namespace, "app=sock-shop", timeout, 2, clientset); err != nil {
		log.Errorf("err: %v", err)
		return
	}
}

//InstallationJaeger is creating and checking status of Jaeger service
func InstallationJaeger(timeout int, clientset *kubernetes.Clientset) {
	if err := CreateNamespace(clientset, "jaeger"); err != nil {
		log.Info("[Status]: Namespace already exist!")
	}

	if err := CreateSockShop("/var/run/jaeger.yaml", "jaeger"); err != nil {
		log.Errorf("Failed to install jaeger, err: %v", err)
		return
	}
	if err := CheckApplicationStatus("jaeger", "app=jaeger", timeout, 2, clientset); err != nil {
		log.Errorf("err: %v", err)
		return
	}
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
