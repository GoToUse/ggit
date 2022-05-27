package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/olekukonko/tablewriter"
)

// RunCommand is a command func which can print the output in real time.
func RunCommand(name string, args ...string) error {
	fmt.Println("Command:", append([]string{name}, args...))
	separator := center(strings.ToUpper(args[0]), 40, "*")
	fmt.Println(separator)
	cmd := exec.Command(name, args...)

	stdout, err := cmd.StdoutPipe()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err != nil {
		fmt.Println("Error details:", err)
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(readString)
		}
	}()

	if err = cmd.Start(); err != nil {
		fmt.Println("Error details:", err)
		return err
	}

	if err = cmd.Wait(); err != nil {
		fmt.Println("Error details:", err)
		return err
	}

	wg.Wait()
	return nil
}

// lookGitPath go to find out where the git is.
func lookGitPath() string {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return ""
	}
	return gitPath
}

// getGitFile return the absolute path of git.
func getGitFile() string {
	file := os.Getenv("GIT")
	if file != "" {
		return file
	}
	gitPath := lookGitPath()
	if gitPath == "" {
		return DefaultGitPath
	}
	return gitPath
}

// retrieveHost get the host of originURL
func retrieveHost(originURL string) string {
	URL, err := url.Parse(originURL)
	if err != nil {
		panic(err)
	}

	return URL.Host
}

// sortHost sorted hosts by RTT value.
func sortHost(originURLList []string) SortHost {
	separator := center("\U0001F973 Sort By Ping RTT Value \U0001F973", 80, "#")
	fmt.Println(separator)

	var rttMapList SortHost
	for _, v := range originURLList {
		wg1.Add(1)
		go func(v string) {
			defer wg1.Done()
			host := retrieveHost(v)
			addr, err := net.LookupCNAME(host)
			if err != nil {
				// Terminates this goroutine
				runtime.Goexit()
			}

			pinger, err := ping.NewPinger(addr)
			if err != nil {
				log.Printf("ping.NewPinger err: %v", err)
				// Terminates this goroutine
				runtime.Goexit()
			}

			fmt.Printf("PING %s (%s)\n", pinger.Addr(), pinger.IPAddr())
			pinger.Count = 5
			pinger.Interval = 500 * time.Millisecond
			pinger.Timeout = 2 * time.Second
			err = pinger.Run()
			stats := pinger.Statistics()
			if err != nil {
				log.Fatalf("pinger.Run err: %v", err)
			}
			rttMapList = append(rttMapList, rttHost{hostName: v, avgRtt: stats.AvgRtt})
			fmt.Printf("%s done!\n", pinger.Addr())
		}(v)
	}
	wg1.Wait()
	sort.SliceStable(rttMapList, func(i, j int) bool {
		return rttMapList[i].avgRtt < rttMapList[j].avgRtt
	})
	return rttMapList
}

// ToMapSetE converts a slice or array to map[any]struct{} with error
func ToMapSetE(i any) (map[any]struct{}, error) {
	// judge the validation of the input
	if i == nil {
		return nil, fmt.Errorf("unable to converts %#v of type %T to map[interface{}]struct{}", i, i)
	}

	kind := reflect.TypeOf(i).Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, fmt.Errorf("the input %#v of type %T isn't a slice or array", i, i)
	}

	// execute the convert
	v := reflect.ValueOf(i)
	vLength := v.Len()
	m := make(map[any]struct{}, vLength)
	for i := 0; i < vLength; i++ {
		m[v.Index(i).Interface()] = struct{}{}
	}

	return m, nil
}

// RenderTable 将host、rtt数据以表格形式呈现
func RenderTable(data []rttHost) {
	var newData [][]string
	for _, item := range data {
		host, rtt := item.hostName, item.avgRtt
		newData = append(newData, []string{host, rtt.String()})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Host", "Round-Trip Time"})

	for _, v := range newData {
		table.Append(v)
	}
	table.Render() // Send output
}

// ObjectsAreEqual determines if two objects are considered equal.
//
// This function does no assertion of any kind.
func ObjectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

// Keys return all keys of map.
func Keys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// HostValues return all values of MirrorUrls.
func HostValues(m map[string][]string) []string {
	values := make([]string, 0, 10)

	for _, v := range m {
		for _, vv := range v {
			values = append(values, vv)
		}
	}

	return values
}

// FindKey find the key by explicit value in Map.
func FindKey(in map[string][]string, predicate any) string {
	for key, value := range in {
		if getType(predicate) == "[]string" {
			if ObjectsAreEqual(value, predicate) {
				return key
			}
		} else if getType(predicate) == "string" {
			valueM, err := ToMapSetE(value)

			if err != nil {
				panic(fmt.Sprintf("FindKey error: %v", err))
			}

			if _, ok := valueM[predicate]; ok {
				return key
			}
		}
	}

	return ""
}

func getType(in any) string {
	return reflect.TypeOf(in).String()
}

func getElemType(in any) string {
	return reflect.TypeOf(in).Elem().String()
}

// center like `str.center` function in python.
func center(s string, n int, fill string) string {
	sLen := len(s)
	div := (n - sLen) / 2
	return strings.Repeat(fill, div) + fmt.Sprintf(" %s ", s) + strings.Repeat(fill, div)
}

// insertElem insert string element into any index of an array.
func insertElem(oriArr []string, position int, elem string) []string {
	oriArr = append(oriArr, "null")
	copy(oriArr[position+1:], oriArr[position:])
	oriArr[position] = elem
	return oriArr
}

// Retry can try to re-run the task if it occurred some temp errors.
func Retry(tryTimes int, sleep time.Duration, callback CallBack) error {
	tipStr := fmt.Sprintf("✨✨✨ Will attempt to retry %d times️ ✨✨✨", tryTimes)
	seperator := center(tipStr, 80, "#")
	fmt.Println(seperator)
	for i := 1; i <= tryTimes; i++ {
		err := callback()
		if err == nil {
			return nil
		}

		if i == tryTimes {
			fmt.Println(fmt.Sprintf("Warning: You have reached the maximum attempts, see error info [%s]", err.Error()))
			fmt.Println(center("💥💥💥I'm a delimiter💥💥💥", 80, "#"))
			return err
		}
		time.Sleep(sleep)
	}
	return nil
}
