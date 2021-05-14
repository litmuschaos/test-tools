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
)

//QPSVars will carry all the params for functionality of end point
type QPSVars struct {
	timeSumsMu    sync.RWMutex
	totalQueries  int
	qpsValue      string
	totalReqCount string
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
	}
	return &qpsVars
}

// It will count the mean value of total queries available in given time interval and update
// end point with query per second(qpsValue).
func runDataLoop(vars *QPSVars) {
	queue := list.New()

	timeInterval, _ := strconv.Atoi(os.Getenv("TIME")) //The time period in second to calculate mean value of QPS.
	url := os.Getenv("URL")                            //URL of endpoint metics
	route := os.Getenv("ROUTE")                        //route is a endpoint of application to get qpsValue
	flag.Parse()

	// Within an infinite loop, lock the mutex and
	// increment our value, then sleep for 1 second until
	// the next time we need to get a value.
	start := time.Now()
	for {

		req, err := GetRequests(url, route, vars)
		if err != nil {
			log.Errorf("err: %v", err)
			vars.qpsValue = strconv.Itoa(0)
			time.Sleep(1 * time.Second)
			continue
		}

		second := int(time.Now().Sub(start).Seconds())
		reqs, _ := strconv.Atoi(req)
		vars.totalQueries = reqs

		if second <= timeInterval {
			queue.PushBack(reqs)
			vars.qpsValue = strconv.Itoa(int(math.Abs(float64(100 * (reqs / queue.Len())))))
		} else {
			front := queue.Front()
			queue.Remove(front)
			queue.PushBack(reqs)
			vars.totalQueries -= front.Value.(int)
			vars.qpsValue = strconv.Itoa(int(math.Abs(float64(100 * (vars.totalQueries / timeInterval)))))
		}
		log.Infof("[Status]: Current total requests : ", req)
		log.Infof("[Status]: Current QPS value is   : ", vars.qpsValue)
		time.Sleep(1 * time.Second)
	}
}

//GetRequests will fetch the response from metrics and calculate the total requests from front-end of sock-shop.
func GetRequests(url string, route string, vars *QPSVars) (string, error) {

	response, err := http.Get(url)
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
		if strings.Contains(string(metricsSplited[i]), `request_duration_seconds_count{service="front-end",method="get",route="`+route+`",status_code="200`) {
			metricsValue := strings.Split(metricsSplited[i], " ")
			vars.totalReqCount = metricsValue[1]
			return vars.totalReqCount, nil
		}
	}

	return vars.totalReqCount, fmt.Errorf("Provided route is incorrect")
}
