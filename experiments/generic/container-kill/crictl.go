package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/litmuschaos/test-tools/pkg/environment"
	"github.com/litmuschaos/test-tools/pkg/events"
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
	eventsDetails := types.EventDetails{}

	//Getting kubeConfig and Generate ClientSets
	if err := clients.GenerateClientSetFromKubeConfig(); err != nil {
		log.Fatalf("Unable to Get the kubeconfig due to %v", err)
	}

	//Fetching all the ENV passed for the runner pod
	log.Info("[PreReq]: Getting the ENV variables")
	environment.GetENV(&experimentsDetails, "container-kill")

	//Obtain the pod ID of the application pod
	podID, err := GetPodID(&experimentsDetails)
	if err != nil {
		log.Fatalf("Unable to get the pod id %v", err)
	}

	log.Infof("Killing the containers of pod with PodID: %v", podID)

	err = KillContainer(&experimentsDetails, clients, podID, &eventsDetails)
	if err != nil {
		log.Fatalf("container-kill chaos terminated due to %v", err)
	}

}

// KillContainer kills the application container
func KillContainer(experimentsDetails *types.ExperimentDetails, clients environment.ClientSets, podID string, eventsDetails *types.EventDetails) error {

	ChaosStartTimeStamp := time.Now().Unix()

	for iteration := 0; iteration < experimentsDetails.Iterations; iteration++ {

		//Obtain the container ID through Pod
		containerID, err := GetContainerID(experimentsDetails, podID)
		if err != nil {
			return errors.Errorf("Unable to get the container id, %v", err)
		}
		log.Infof("Killing the container with containerID: %v", containerID)

		if experimentsDetails.EngineName != "" {
			msg := "Injecting " + experimentsDetails.ExperimentName + " chaos on application pod"
			environment.SetEventAttributes(eventsDetails, types.ChaosInject, msg)
			events.GenerateEvents(experimentsDetails, clients, eventsDetails)
		}

		// killing the application container
		StopContainer(containerID)

		//Waiting for the chaos interval after chaos injection
		if experimentsDetails.ChaosInterval != 0 {
			log.Infof("[Wait]: Wait for the chaos interval %vs", strconv.Itoa(experimentsDetails.ChaosInterval))
			waitForChaosInterval(experimentsDetails)
		}

		//Check the status of restarted container
		err = CheckContainerStatus(experimentsDetails, clients)
		if err != nil {
			return errors.Errorf("Application container is not running, %v", err)
		}

		ChaosCurrentTimeStamp := time.Now().Unix()
		chaosDiffTimeStamp := ChaosCurrentTimeStamp - ChaosStartTimeStamp

		// terminating the execution after the timestamp exceed the total chaos duration
		if int(chaosDiffTimeStamp) >= experimentsDetails.ChaosDuration {
			break
		}

	}
	log.Infof("[Completion]: %v chaos is done", experimentsDetails.ExperimentName)
	return nil

}

//GetPodID derive the pod-id of the application pod
func GetPodID(experimentsDetails *types.ExperimentDetails) (string, error) {

	cmd := exec.Command("crictl", "pods")
	stdout, _ := cmd.Output()

	pods := removeExtraSpaces(stdout)
	for i := 0; i < len(pods)-1; i++ {
		attributes := strings.Split(pods[i], " ")
		// fmt.Printf("podlist: %v", attributes)
		if attributes[3] == experimentsDetails.ApplicationPod {
			return attributes[0], nil
		}

	}

	return "", fmt.Errorf("The application pod is unavailable")
}

//GetContainerID  derive the container-id of the application container
func GetContainerID(experimentsDetails *types.ExperimentDetails, podID string) (string, error) {

	cmd := exec.Command("crictl", "ps")
	stdout, _ := cmd.Output()
	containers := removeExtraSpaces(stdout)

	for i := 0; i < len(containers)-1; i++ {
		attributes := strings.Split(containers[i], " ")
		if attributes[4] == experimentsDetails.ApplicationContainer && attributes[6] == podID {
			return attributes[0], nil
		}

	}

	return "", fmt.Errorf("The application container is unavailable")

}

//StopContainer kill the application container
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

//waitForChaosInterval waits for the given ramp time duration (in seconds)
func waitForChaosInterval(experimentsDetails *types.ExperimentDetails) {
	time.Sleep(time.Duration(experimentsDetails.ChaosInterval) * time.Second)
}

// removeExtraSpaces remove all the extra spaces present in output of crictl commands
func removeExtraSpaces(arr []byte) []string {
	bytesSlice := make([]byte, len(arr))
	index := 0
	count := 0
	for i := 0; i < len(arr); i++ {
		count = 0
		for arr[i] == 32 {
			count++
			i++
			if i >= len(arr) {
				break
			}
		}
		if count > 1 {
			bytesSlice[index] = 32
			index++
		}
		bytesSlice[index] = arr[i]
		index++

	}
	return strings.Split(string(bytesSlice), "\n")
}
