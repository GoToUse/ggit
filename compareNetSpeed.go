package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

var wg1 sync.WaitGroup

func retrieveHost(originURL string) string {
	orgURLList := strings.Split(originURL, "//")
	host := orgURLList[1]
	return strings.TrimSuffix(host, "/")
}

func main() {
	rttMap := make(map[string]interface{})
	var rttMapList []struct {
		hostName string
		avgRtt   time.Duration
	}
	for _, v := range DEFAULT_MIRROR_URL_ARRAY {
		wg1.Add(1)
		go func(v string) {
			defer wg1.Done()
			host := retrieveHost(v)
			pinger, err := ping.NewPinger(host)
			if err != nil {
				panic(err)
			}

			fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
			pinger.Count = 5
			pinger.Interval = 500 * time.Microsecond
			pinger.Timeout = 2 * time.Second
			fmt.Printf("%s run...\n", pinger.Addr())
			err = pinger.Run()
			stats := pinger.Statistics()
			rttMap[pinger.Addr()] = stats.AvgRtt
			if err != nil {
				panic(err)
			}
			rttMapList = append(rttMapList, struct {
				hostName string
				avgRtt   time.Duration
			}{hostName: pinger.Addr(), avgRtt: stats.AvgRtt})
			fmt.Printf("%s done!\n", pinger.Addr())
		}(v)
	}
	wg1.Wait()
	fmt.Println(rttMap)
	sort.SliceStable(rttMapList, func(i, j int) bool {
		return rttMapList[i].avgRtt < rttMapList[j].avgRtt
	})
	fmt.Println(rttMapList)
}
