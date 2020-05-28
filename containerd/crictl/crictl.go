package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/litmuschaos/test-tools/pkg/environment"
	"github.com/litmuschaos/test-tools/pkg/log"
	"github.com/litmuschaos/test-tools/pkg/types"
	"github.com/openebs/maya/pkg/util/retry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableSorting:         true,
		DisableLevelTruncation: true,
	})
}

func main() {

	experimentsDetails := types.ExperimentDetails{}
	clients := environment.ClientSets{}

	//Getting kubeConfig and Generate ClientSets
	if err := clients.GenerateClientSetFromKubeConfig(); err != nil {
		log.Fatalf("Unable to Get the kubeconfig due to %v", err)
	}

	//Fetching all the ENV passed for the runner pod
	log.Info("[PreReq]: Getting the ENV variables")
	environment.GetENV(&experimentsDetails)

	//Deriving the chaos iterations
	experimentsDetails.Iterations = GetChaosIteration(&experimentsDetails)
	fmt.Printf("iteration %v\n", experimentsDetails.Iterations)

	//Obtain the pod ID through Pod name
	podID, err := GetPodID(&experimentsDetails)
	if err != nil {
		log.Errorf("Error %v", err)
	}
	log.Infof("PodID %v", podID)

	KillContainer(&experimentsDetails, clients, podID)

}

// KillContainer ..
func KillContainer(experimentsDetails *types.ExperimentDetails, clients environment.ClientSets, podID string) error {

	ChaosStartTimeStamp := time.Now().Unix()

	for x := 0; x < experimentsDetails.Iterations; x++ {

		//Obtain the container ID through Pod name
		containerID, err := GetContainerID(experimentsDetails, podID)
		if err != nil {
			return err
		}
		log.Infof("containerID %v", containerID)

		// stop container
		StopContainer(containerID)

		//Check if the new container is running
		err = CheckContainerStatus(experimentsDetails, clients)
		if err != nil {
			return err
		}

		ChaosCurrentTimeStamp := time.Now().Unix()
		chaosDiffTimeStamp := ChaosCurrentTimeStamp - ChaosStartTimeStamp

		if int(chaosDiffTimeStamp) >= experimentsDetails.ChaosDuration {
			break
		}

	}
	return nil

}

//GetChaosIteration ...
func GetChaosIteration(experimentsDetails *types.ExperimentDetails) int {
	return (experimentsDetails.ChaosDuration / experimentsDetails.ChaosInterval)

}

//GetPodID ...
func GetPodID(experimentsDetails *types.ExperimentDetails) (string, error) {

	cmd := exec.Command("crictl", "pods")
	stdout, _ := cmd.Output()
	ans := string(stdout)
	fmt.Print(string(stdout))

	res1 := strings.Split(ans, "\n")

	for i := 0; i < len(res1)-1; i++ {
		// fmt.Println(res1[i])
		res2 := strings.Split(res1[i], " ")
		if res2[2] == experimentsDetails.ApplicationPod {
			return res2[1], nil
		}

	}

	return "", fmt.Errorf("The application pod is unavailable")

}

//GetContainerID ...
func GetContainerID(experimentsDetails *types.ExperimentDetails, podID string) (string, error) {

	cmd := exec.Command("crictl", "ps")
	stdout, _ := cmd.Output()
	ans := string(stdout)
	fmt.Print(string(stdout))

	res1 := strings.Split(ans, "\n")

	for i := 0; i < len(res1)-1; i++ {
		// fmt.Println(res1[i])
		res2 := strings.Split(res1[i], " ")
		if res2[2] == experimentsDetails.ApplicationContainer && res2[3] == podID {
			return res2[1], nil
		}

	}

	return "", fmt.Errorf("The application container is unavailable")

}

//StopContainer ...
func StopContainer(containerID string) {

	cmd := exec.Command("crictl", "stop", string(containerID))
	stdout, _ := cmd.Output()
	fmt.Print(string(stdout))

}

// CheckContainerStatus checks the status of the application container
func CheckContainerStatus(experimentsDetails *types.ExperimentDetails, clients environment.ClientSets) error {
	err := retry.
		Times(90).
		Wait(2 * time.Second).
		Try(func(attempt uint) error {
			pod, err := clients.KubeClient.CoreV1().Pods(experimentsDetails.AppNS).Get(experimentsDetails.ApplicationPod, v1.GetOptions{})
			if err != nil {
				return errors.Errorf("Unable to get the pod, err: %v", err)
			}
			err = nil
			for _, container := range pod.Status.ContainerStatuses {
				if container.Ready != true {
					return errors.Errorf("containers are not yet in running state")
				}
				log.InfoWithValues("The running status of container are as follows", logrus.Fields{
					"container": container.Name, "Pod": pod.Name, "Status": pod.Status.Phase})
			}

			return nil
		})
	if err != nil {
		return err
	}
	return nil
}
