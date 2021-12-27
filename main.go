package main

import (
	"bufio"
	"fmt"
	"github.com/go-ping/ping"
	"io"
	"log"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	DEFAULT_GIT_PATH      string = "/usr/local/bin/git"
	DEFAULT_GITHUB_URL    string = "https://github.com/"
	DEFAULT_GITHUB_SUFFIX string = ".git"
)

type (
	rttHost struct {
		hostName string
		avgRtt   time.Duration
	}
	Args []string
	SortHost []rttHost
)

var (
	wg                       sync.WaitGroup
	wg1                      sync.WaitGroup
	DEFAULT_MIRROR_URL_ARRAY = []string{
		"https://hub.fastgit.org/",
		"https://github.com.cnpmjs.org/",
		//"https://gitclone.com/",
		"https://github.wuyanzheshui.workers.dev/",
	}
)

func RunCommand(name string, args ...string) error {
	fmt.Println("Command: ", append([]string{name}, args...))
	seperator := center(strings.ToUpper(args[0]), 60, "*")
	fmt.Println(seperator)
	cmd := exec.Command(name, args...)

	stdout, err := cmd.StdoutPipe()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	defer stdout.Close()

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

func lookGitPath() string {
	gitPath, err := exec.LookPath("git")
	if err != nil {
		return ""
	}
	return gitPath
}

func getGitFile() string {
	file := os.Getenv("GIT")
	if file != "" {
		return file
	}
	gitPath := lookGitPath()
	if gitPath == "" {
		return DEFAULT_GIT_PATH
	}
	return gitPath
}

func ggitClone(args Args, mirrorUrl string) error {
	var oldUrl, newUrl, folderName string

	if strings.HasPrefix(args[2], DEFAULT_GITHUB_URL) {
		oldUrl = args[2]
		// 特别处理
		u, err := url.Parse(mirrorUrl)
		if err != nil {
			log.Panicf("%s is wrong, see details[%s]", mirrorUrl, err.Error())
		}
		if strings.Contains(mirrorUrl, "https://gitclone.com/") {
			ref, _ := u.Parse("github.com")
			githubCloneUrl := fmt.Sprintf("%s/", ref)
			newUrl = strings.ReplaceAll(oldUrl, DEFAULT_GITHUB_URL, githubCloneUrl)
		} else {
			newUrl = strings.ReplaceAll(oldUrl, DEFAULT_GITHUB_URL, mirrorUrl)
		}
		args[2] = newUrl
		folderNameArr := strings.Split(oldUrl, "/")
		folderName = folderNameArr[len(folderNameArr)-1]
		if strings.HasSuffix(folderName, ".git") {
			folderName = strings.Split(folderName, ".git")[0]
		}
		fmt.Println("Folder name:", folderName)
	} else {
		fmt.Printf("DEBUG: args[2]: %s\n", args[2])
		log.Fatal("github仓库地址有误, 请检查是否符合 [https://github.com/xxx/xxx.git] 标准路径.")
	}

	args[0] = getGitFile()
	err := RunCommand(args[0], args[1:]...)
	if err != nil || len(newUrl) == 0 || len(folderName) == 0 {
		retryErr := Retry(3, 3 * time.Second, func() error {
			fErr := RunCommand(args[0], args[1:]...)
			return fErr
		})
		if retryErr != nil {
			// 如果当前url不能正常工作，那么初始化args[2]的值
			args[2] = oldUrl
			return err
		}
	}
	fmt.Println("Clone done!!!")

	err = os.Chdir(folderName)
	if err != nil {
		panic(err)
	}

	restoreCmd := "remote set-url origin " + oldUrl
	err = RunCommand(args[0], strings.Fields(restoreCmd)...)
	if err != nil {
		retryErr := Retry(3, 3 * time.Second, func() error {
			fErr := RunCommand(args[0], args[1:]...)
			return fErr
		})
		if retryErr != nil {
			// 如果当前url不能正常工作，那么初始化args[2]的值
			args[2] = oldUrl
			// TODO: if error, delete this folder.
			panic(err)
		}
	}

	fmt.Println("Set remote done!!!")
	return nil
}

func retrieveHost(originURL string) string {
	orgURLList := strings.Split(originURL, "//")
	host := orgURLList[1]
	return strings.TrimSuffix(host, "/")
}

func sortHost(originURLList []string) SortHost {
	seperator := center("Sort by ping rtt value", 60, "#")
	fmt.Println(seperator)
	var rttMapList SortHost
	for _, v := range originURLList {
		wg1.Add(1)
		go func(v string) {
			defer wg1.Done()
			host := retrieveHost(v)
			pinger, err := ping.NewPinger(host)
			if err != nil {
				panic(err)
			}

			fmt.Printf("PING %s (%s)\n", pinger.Addr(), pinger.IPAddr())
			pinger.Count = 5
			pinger.Interval = 500 * time.Millisecond
			pinger.Timeout = 2 * time.Second
			err = pinger.Run()
			stats := pinger.Statistics()
			if err != nil {
				panic(err)
			}
			rttMapList = append(rttMapList, struct {
				hostName string
				avgRtt   time.Duration
			}{hostName: v, avgRtt: stats.AvgRtt})
			fmt.Printf("%s done!\n", pinger.Addr())
		}(v)
	}
	wg1.Wait()
	sort.SliceStable(rttMapList, func(i, j int) bool {
		return rttMapList[i].avgRtt < rttMapList[j].avgRtt
	})
	fmt.Printf("Sorted list: %v\n", rttMapList)
	return rttMapList
}

func GgitClone(args Args) {
	var initTimes int
	sortHost := sortHost(DEFAULT_MIRROR_URL_ARRAY)
	for _, v := range sortHost {
		mirrorUrl := v.hostName
		fmt.Println("# Current mirror's url is: ", mirrorUrl)
		err := ggitClone(args, mirrorUrl)
		if err != nil {
			initTimes += 1
			continue
		}
		fmt.Println("All done!!!")
		return
	}

	if initTimes == len(DEFAULT_MIRROR_URL_ARRAY) {
		log.Fatal("Sorry: All mirrors are unusable.")
	}
}

type CallBack func() error

// Retry can try to re-run the task if it occurred some temp errors.
func Retry(tryTimes int, sleep time.Duration, callback CallBack) error {
	tipStr := fmt.Sprintf("Will attempt to retry %d times", tryTimes)
	seperator := center(tipStr, 60, "#")
	fmt.Println(seperator)
	for i := 1; i <= tryTimes; i++ {
		err := callback()
		if err == nil {
			return nil
		}

		if i == tryTimes {
			fmt.Println(fmt.Sprintf("Warning: You have reached the maximum attempts, see error info [%s]", err.Error()))
			return err
		}
		time.Sleep(sleep)
	}
	return nil
}

// center like `str.center` function in python.
func center(s string, n int, fill string) string {
	div := n / 2
	return strings.Repeat(fill, div) + s + strings.Repeat(fill, div)
}

func main() {
	cmdArgs := os.Args
	fmt.Println(cmdArgs)
	if len(cmdArgs) > 2 &&
		cmdArgs[1] == "clone" &&
		strings.HasSuffix(cmdArgs[2], DEFAULT_GITHUB_SUFFIX) {
		GgitClone(cmdArgs)
	} else if len(cmdArgs) > 3 {
		log.Fatal("请输入[ggit clone https://github.com/xxx/xxx.git]这样的命令格式")
	} else {
		log.Fatal("请输入[ggit clone https://github.com/xxx/xxx.git]这样的命令格式")
	}
}
