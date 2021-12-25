package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	DEFAULT_GIT_PATH      string = "/usr/local/bin/git"
	DEFAULT_GITHUB_URL    string = "https://github.com/"
	DEFAULT_GITHUB_SUFFIX string = ".git"
)

var (
	wg                       sync.WaitGroup
	SEPERATOR                = strings.Repeat("*", 30)
	DEFAULT_MIRROR_URL_ARRAY = [...]string{
		"https://github.com.cnpmjs.org/",
		"https://hub.fastgit.org/",
		"https://gitclone.com/",
		"https://github.wuyanzheshui.workers.dev/",
	}
)

type Args []string

func RunCommand(name string, arg ...string) error {
	fmt.Println("Command: ", append([]string{name}, arg...))
	fmt.Printf("%s %s %s\n", SEPERATOR, strings.ToUpper(arg[0]), SEPERATOR)
	cmd := exec.Command(name, arg...)

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

	for i := 2; i < len(args); i++ {
		if strings.HasPrefix(args[i], DEFAULT_GITHUB_URL) {
			oldUrl = args[i]
			// 特别处理
			if strings.Contains(mirrorUrl, "https://gitclone.com/") {
				githubCloneUrl := path.Join(mirrorUrl, "github.com") + "/"
				newUrl = strings.ReplaceAll(oldUrl, DEFAULT_GITHUB_URL, githubCloneUrl)
			} else {
				newUrl = strings.ReplaceAll(oldUrl, DEFAULT_GITHUB_URL, mirrorUrl)
			}
			args[i] = newUrl
			folderNameArr := strings.Split(oldUrl, "/")
			folderName = folderNameArr[len(folderNameArr)-1]
			if strings.HasSuffix(folderName, ".git") {
				folderName = strings.Split(folderName, ".git")[0]
			}
			fmt.Println("Folder name:", folderName)
		} else {
			log.Fatal("github仓库地址有误, 请检查是否符合 [https://github.com/xxx/xxx.git] 标准路径.")
		}
	}

	args[0] = getGitFile()
	err := RunCommand(args[0], args[1:]...)
	if err != nil || len(newUrl) == 0 || len(folderName) == 0 {
		retryErr := Retry(3, time.Second, func() error {
			fErr := RunCommand(args[0], args[1:]...)
			return fErr
		})
		if retryErr != nil {
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
		retryErr := Retry(3, time.Second, func() error {
			fErr := RunCommand(args[0], args[1:]...)
			return fErr
		})
		if retryErr != nil {
			panic(err)
		}
	}

	fmt.Println("Set remote done!!!")
	return nil
}

func GgitClone(args Args) {
	var initTimes int
	for _, mirrorUrl := range DEFAULT_MIRROR_URL_ARRAY {
		fmt.Println("Current mirror's url is: ", mirrorUrl)
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

func Retry(tryTimes int, sleep time.Duration, callback CallBack) error {
	fmt.Printf("Will attempt to retry %d times.\n", tryTimes)
	for i := 1; i <= tryTimes; i++ {
		err := callback()
		if err == nil {
			return nil
		}

		if i == tryTimes {
			panic(fmt.Sprintf("You have reached the maximum attempts, see error info [%s]", err.Error()))
			return err
		}
		time.Sleep(sleep)
	}
	return nil
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
