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
	application   string
}

func main() {
	log.Info("[Status]: Starting QPS provider...")

	qpsVars := initialiseVars()
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
		application:   os.Getenv("APPLICATION_NAME"),
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
	sumQuery := 0
	skipFirst := true
	vars.timeInterval, _ = strconv.Atoi(os.Getenv("TIME")) //The time period in second to calculate mean value of QPS.
	flag.Parse()

	// Within an infinite loop, lock the mutex and
	// increment our value, then sleep for 1 second until
	// the next time we need to get a value.
	start := time.Now()
	for {
		vars.urlList, err = getURL(vars, clientset)
		if err != nil {
			log.Errorf("Unable to find pods in namespace, err: %v", err)
			time.Sleep(1 * time.Second)
			vars.totalReqCount = "0"
		}
		req, err := GetRequests(vars.urlList, vars.route, vars)
		if err != nil {
			vars.qpsValue = strconv.Itoa(0)
			time.Sleep(1 * time.Second)
		}

		if skipFirst {
			sumQuery = 0
			skipFirst = false
			vars.totalQueries, _ = strconv.Atoi(req)
			vars.totalReqCount = "0"
			continue
		}
		second := int(time.Now().Sub(start).Seconds())
		reqs, _ := strconv.Atoi(req)
		vars.totalReqCount = "0"

		if second <= vars.timeInterval {
			if reqs-vars.totalQueries > 0 {
				queue.PushBack(reqs - vars.totalQueries)
				sumQuery += reqs - vars.totalQueries
				vars.qpsValue = strconv.Itoa(int(math.Abs(float64(100 * int((sumQuery)) / vars.timeInterval))))
			} else {
				queue.PushBack(0)
				vars.qpsValue = strconv.Itoa(int(math.Abs(float64(100 * int((sumQuery)) / vars.timeInterval))))
			}
		} else {
			front := queue.Front()
			sumQuery -= front.Value.(int)
			queue.Remove(front)
			if reqs-vars.totalQueries > 0 {
				queue.PushBack(reqs - vars.totalQueries)
				sumQuery += reqs - vars.totalQueries
			} else {
				queue.PushBack(0)
			}
			vars.qpsValue = strconv.Itoa(int(math.Abs(float64(100 * (sumQuery / vars.timeInterval)))))
		}
		vars.totalQueries = reqs
		log.Infof("[Status]: Current Total Requests : ", req)
		log.Infof("[Status]: Current Query Rate : ", vars.qpsValue)
		time.Sleep(1 * time.Second)
	}
}

// getURL will list the IPs for all the pods exporting metrics
func getURL(vars *QPSVars, clientset *kubernetes.Clientset) ([]string, error) {
	vars.urlList = []string{}

	switch strings.ToLower(vars.application) {
	case "postgres":
		podSpec, err := clientset.CoreV1().Pods(vars.namespace).List(metav1.ListOptions{LabelSelector: vars.appLabel})
		if err != nil && len(podSpec.Items) == 0 {
			return []string{}, err
		}
		for _, pod := range podSpec.Items {
			if strings.Contains(string(pod.ObjectMeta.Annotations["status"]), "master") {
				vars.urlList = append(vars.urlList, string(`http://`+strings.Replace(pod.Status.PodIP, ".", "-", -1)+`.`+vars.route))
				break
			}
		}
		if len(vars.urlList) == 0 {
			return []string{}, errors.Errorf("Unable to find the pods with master role in namespace")
		}
	case "sock-shop":
		vars.urlList = append(vars.urlList, string(vars.route))
	default:
		return []string{}, errors.Errorf("Application '%v' not supported in app-qps-test", vars.application)
	}
	return vars.urlList, nil
}

//GetRequests will fetch the response from metrics and calculate the total requests from front-end of sock-shop.
func GetRequests(urlList []string, route string, vars *QPSVars) (string, error) {

	for index := range urlList {
		response, err := http.Get(urlList[index])
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
