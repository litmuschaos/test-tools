package main

import (
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/litmuschaos/test-tools/pkg/log"
)

var (
	timeSumsMu sync.RWMutex
	timeSums   int64
	qps        string
	prevReq    int
	sum        int
)

func main() {
	log.Info("[Status]: Starting QPS provider...")
	go runDataLoop()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		timeSumsMu.RLock()
		defer timeSumsMu.RUnlock()
		fmt.Fprint(w, qps)
	})
	http.ListenAndServe(":8080", nil)
}

//// Start the goroutine that will sum the current time
// once per second.
// Create a handler that will read-lock the mutext and
// write the summed time to the client

func runDataLoop() {
	queue := list.New()
	timeInt := os.Getenv("TIME") //The number of seconds to calculate mean value of QPS.
	url := os.Getenv("URI")      // URL of endpoint
	flag.Parse()

	timeOfInt, _ := strconv.Atoi(timeInt)

	for {
		// Within an infinite loop, lock the mutex and
		// increment our value, then sleep for 1 second until
		// the next time we need to get a value.
		start := time.Now()
		timeSumsMu.Lock()
		timeSums += time.Now().Unix()
		timeSumsMu.Unlock()
		prevReq = 0
		for i := 1; ; i++ {
			req, err := GetRequests(url)

			fmt.Println(req)
			if err != nil {
				qps = strconv.Itoa(0)
				fmt.Printf("%s", err)
				return
			}
			now := time.Now()
			diff := now.Sub(start)
			second := int(diff.Seconds())
			reqs, err := strconv.Atoi(req)
			qps = string(strconv.Itoa(reqs - prevReq))
			prevReq = reqs

			if second < 20 {
				sum = reqs
				queue.PushBack(reqs)
				//	fmt.Println("QPS :%d Sec : %d", qps, second)
			} else {
				front := queue.Front()
				queue.Remove(front)
				queue.PushBack(reqs)
				sum = reqs
				sum -= front.Value.(int)
				qps = string(sum / timeOfInt)
				fmt.Println("QPS :", sum/timeOfInt)
				fmt.Println("Sec :", second)
			}
			log.Infof("[Status]: Current total requests : ", req)
			log.Infof("[Status]: Current QPS value is : ", qps)
		}
	}
}

//GetRequests will fetch the responce from metrics and calculate the total requests from front-end of sock-shop.
func GetRequests(url string) (string, error) {
	var totalCount string
	time.Sleep(1 * time.Second)
	//"http://front-end.sock-shop.svc.cluster.local/metrics"
	response, err := http.Get(url)
	if err != nil {
		qps = strconv.Itoa(0)
		log.Errorf("Failed to fetch responce, err: %v", err)
		return "", err
	} else {
		defer response.Body.Close()
		metric, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return "", err
		}

		metrics := string(metric)
		count := 92
		it := strings.Index(metrics, "request_duration_seconds_count")

		var str string
		for i := it + count + 1; ; i++ {
			if string(metrics[i]) >= "0" && string(metrics[i]) <= "9" {
				str += string(metrics[i])
			} else {
				break
			}
		}
		totalCount = str
	}
	return totalCount, nil
}
