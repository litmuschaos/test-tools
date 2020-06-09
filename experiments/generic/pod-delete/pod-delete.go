package main

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/litmuschaos/test-tools/pkg/environment"
	"github.com/litmuschaos/test-tools/pkg/events"
	"github.com/litmuschaos/test-tools/pkg/log"
	"github.com/litmuschaos/test-tools/pkg/math"
	"github.com/litmuschaos/test-tools/pkg/status"
	"github.com/litmuschaos/test-tools/pkg/types"
	"github.com/openebs/maya/pkg/util/retry"
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

	var err error
	experimentsDetails := types.ExperimentDetails{}
	clients := environment.ClientSets{}
	eventsDetails := types.EventDetails{}

	//Getting kubeConfig and Generate ClientSets
	if err := clients.GenerateClientSetFromKubeConfig(); err != nil {
		log.Fatalf("Unable to Get the kubeconfig due to %v", err)
	}

	//Fetching all the ENV passed for the runner pod
	log.Info("[PreReq]: Getting the ENV variables")
	environment.GetENV(&experimentsDetails, "pod-delete")

	err = PodDeleteChaos(&experimentsDetails, clients, &eventsDetails)

	if err != nil {
		log.Fatalf("Unable to delete the application pods, due to %v", err)
	}

}

//PodDeleteChaos deletes the random single/multiple pods
func PodDeleteChaos(experimentsDetails *types.ExperimentDetails, clients environment.ClientSets, eventsDetails *types.EventDetails) error {

	//ChaosStartTimeStamp contains the start timestamp, when the chaos injection begin
	ChaosStartTimeStamp := time.Now().Unix()
	var GracePeriod int64 = 0

	for x := 0; x < experimentsDetails.Iterations; x++ {
		//Getting the list of all the target pod for deletion
		targetPodList, err := PreparePodList(experimentsDetails, clients)
		if err != nil {
			return err
		}
		log.InfoWithValues("[Info]: Killing the following pods", logrus.Fields{
			"PodList": targetPodList})

		if experimentsDetails.EngineName != "" {
			msg := "Injecting " + experimentsDetails.ExperimentName + " chaos on application pod"
			environment.SetEventAttributes(eventsDetails, types.ChaosInject, msg)
			events.GenerateEvents(experimentsDetails, clients, eventsDetails)
		}

		//Deleting the application pod
		for _, pods := range targetPodList {
			if experimentsDetails.Force == true {
				err = clients.KubeClient.CoreV1().Pods(experimentsDetails.AppNS).Delete(pods, &v1.DeleteOptions{GracePeriodSeconds: &GracePeriod})
			} else {
				err = clients.KubeClient.CoreV1().Pods(experimentsDetails.AppNS).Delete(pods, &v1.DeleteOptions{})
			}
		}
		if err != nil {
			return err
		}

		//Waiting for the chaos interval after chaos injection
		if experimentsDetails.ChaosInterval != 0 {
			log.Infof("[Wait]: Wait for the chaos interval %vs", strconv.Itoa(experimentsDetails.ChaosInterval))
			waitForChaosInterval(experimentsDetails)
		}
		//Verify the status of pod after the chaos injection
		log.Info("[Status]: Verification for the recreation of application pod")
		err = status.CheckApplicationStatus(experimentsDetails.AppNS, experimentsDetails.AppLabel, clients)

		//ChaosCurrentTimeStamp contains the current timestamp
		ChaosCurrentTimeStamp := time.Now().Unix()

		//ChaosDiffTimeStamp contains the difference of current timestamp and start timestamp
		//It will helpful to track the total chaos duration
		chaosDiffTimeStamp := ChaosCurrentTimeStamp - ChaosStartTimeStamp

		if int(chaosDiffTimeStamp) >= experimentsDetails.ChaosDuration {
			break
		}

	}
	log.Infof("[Completion]: %v chaos is done", experimentsDetails.ExperimentName)

	return nil
}

//waitForChaosInterval waits for the given ramp time duration (in seconds)
func waitForChaosInterval(experimentsDetails *types.ExperimentDetails) {
	time.Sleep(time.Duration(experimentsDetails.ChaosInterval) * time.Second)
}

//PreparePodList derive the list of target pod for deletion
//It is based on the KillCount value
func PreparePodList(experimentsDetails *types.ExperimentDetails, clients environment.ClientSets) ([]string, error) {

	var targetPodList []string

	err := retry.
		Times(90).
		Wait(2 * time.Second).
		Try(func(attempt uint) error {
			pods, err := clients.KubeClient.CoreV1().Pods(experimentsDetails.AppNS).List(v1.ListOptions{LabelSelector: experimentsDetails.AppLabel})
			if err != nil || len(pods.Items) == 0 {
				return errors.Errorf("Unable to get the pod, err: %v", err)
			}
			index := rand.Intn(len(pods.Items))
			//Adding the first pod only, if KillCount is not set or 0
			//Otherwise derive the min(KIllCount,len(pod_list)) pod
			if experimentsDetails.KillCount == 0 {
				targetPodList = append(targetPodList, pods.Items[index].Name)
			} else {
				for i := 0; i < math.Minimum(experimentsDetails.KillCount, len(pods.Items)); i++ {
					targetPodList = append(targetPodList, pods.Items[index].Name)
					index = (index + 1) % len(pods.Items)
				}
			}

			return nil
		})

	return targetPodList, err
}
