package main

import (
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/litmuschaos/test-tools/pkg/log"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//QPSVars will carry all the params for functionality of end point
type QPSVars struct {
	timeSumsMu    sync.RWMutex
	totalQueries  int
	qpsValue      string
	totalReqCount string
	namespace     string
	appLabel      string
	timeInterval  int
	route         string
	urlList       []string
	query         string
}

func main() {
	log.Info("[Status]: Starting QPS provider...")

	qpsVars := initialiseVars()
	log.Info("[Status]: intialised")
	go runDataLoop(qpsVars)

	// Create a handler that will read-lock the mutext and
	// write the summed time to the end point
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		qpsVars.timeSumsMu.RLock()
		defer qpsVars.timeSumsMu.RUnlock()
		fmt.Fprint(w, qpsVars.qpsValue)
	})
	http.ListenAndServe(":8080", nil)
}

//initialiseVars is initialling required variables for app-test
//totalQueries is keeping count of requests in given time interval
//qpsValue is Queries per second
//totalReqCount is total requests to endpoint
func initialiseVars() (vars *QPSVars) {
	qpsVars := QPSVars{
		totalQueries:  0,
		qpsValue:      "0",
		totalReqCount: "0",
		namespace:     os.Getenv("APP_NAMESPACE"),
		appLabel:      os.Getenv("APP_LABEL"),
		timeInterval:  0,
		route:         os.Getenv("ROUTE"),
		query:         os.Getenv("QUERY"),
		urlList:       []string{},
	}
	return &qpsVars
}

// It will count the mean value of total queries available in given time interval and update
// end point with query per second(qpsValue).
func runDataLoop(vars *QPSVars) {

	config, err := getKubeConfig()
	if err != nil {
		return
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	queue := list.New()
	vars.timeInterval, _ = strconv.Atoi(os.Getenv("TIME")) //The time period in second to calculate mean value of QPS.
	flag.Parse()

	// Within an infinite loop, lock the mutex and
	// increment our value, then sleep for 1 second until
	// the next time we need to get a value.
	start := time.Now()
	for {
		vars.urlList, err = getURL(vars, clientset)
		if err != nil {
			log.Errorf("Unable to find the pods in namespace, err: %v", err)
			continue
		}
		req, err := GetRequests(vars.urlList, vars.route, vars)
		if err != nil {
			vars.qpsValue = strconv.Itoa(0)
			time.Sleep(1 * time.Second)
			continue
		}

		second := int(time.Now().Sub(start).Seconds())
		reqs, _ := strconv.Atoi(req)
		vars.totalReqCount = "0"
		vars.totalQueries = reqs

		if second <= vars.timeInterval {
			queue.PushBack(reqs)
			vars.qpsValue = strconv.Itoa(int(math.Abs(float64(100 * (reqs / queue.Len())))))
		} else {
			front := queue.Front()
			queue.Remove(front)
			queue.PushBack(reqs)
			vars.totalQueries -= front.Value.(int)
			vars.qpsValue = strconv.Itoa(int(math.Abs(float64(100 * (vars.totalQueries / vars.timeInterval)))))
		}
		log.Infof("[Status]: Current total requests : ", req)
		log.Infof("[Status]: Current QPS value is   : ", vars.qpsValue)
		time.Sleep(1 * time.Second)
	}
}

// getURL will list the IPs for all the pods exporting metrics
func getURL(vars *QPSVars, clientset *kubernetes.Clientset) ([]string, error) {
	vars.urlList = []string{}
	podSpec, err := clientset.CoreV1().Pods(vars.namespace).List(metav1.ListOptions{LabelSelector: vars.appLabel})
	if err != nil {
		return []string{}, errors.Errorf("Unable to find the pods in namespace, err: %v", err)
	}
	for _, pod := range podSpec.Items {
		vars.urlList = append(vars.urlList, strings.Replace(pod.Status.PodIP, ".", "-", -1))
	}
	return vars.urlList, nil
}

//GetRequests will fetch the response from metrics and calculate the total requests from front-end of sock-shop.
func GetRequests(urlList []string, route string, vars *QPSVars) (string, error) {

	for index := range urlList {
		response, err := http.Get(string("http://" + urlList[index] + "." + vars.route))
		if err != nil {
			return "", err
		}
		defer response.Body.Close()

		metric, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}

		metrics := string(metric)
		metricsSplited := strings.Split(metrics, "\n")
		for i := 0; i < len(metricsSplited); i++ {
			if strings.Contains(string(metricsSplited[i]), vars.query) {
				metricsValue := strings.Split(metricsSplited[i], " ")
				currentCount, _ := strconv.Atoi(vars.totalReqCount)
				newCount, _ := strconv.Atoi(metricsValue[1])
				totalCount := currentCount + newCount
				vars.totalReqCount = strconv.Itoa(totalCount)
			}
		}
	}
	return vars.totalReqCount, nil
}

// GetKubeConfig function derive the kubeconfig
func getKubeConfig() (*rest.Config, error) {
	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()
	// It uses in-cluster config if kubeconfig path is not specified
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	return config, err
}
