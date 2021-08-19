package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var (
	DEFAULT_GIT_FILE   string = "/usr/bin/git"
	DEFAULT_GITHUB_URL string = "https://github.com/"
	DEFAULT_MIRROR_URL string = "https://github.com.cnpmjs.org/"
	SEPERATOR          string = strings.Repeat("*", 30)
	wg                 sync.WaitGroup
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
		return err
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	wg.Wait()
	return nil
}

func getGitFile() string {
	file := os.Getenv("GIT")
	if file != "" {
		return file
	}
	return DEFAULT_GIT_FILE
}

func GgitClone(args Args) {
	var old_url, new_url, folder_name string

	for i := 2; i < len(args); i++ {
		if strings.HasPrefix(args[i], DEFAULT_GITHUB_URL) {
			old_url = args[i]
			new_url = strings.ReplaceAll(old_url, DEFAULT_GITHUB_URL, DEFAULT_MIRROR_URL)
			args[i] = new_url
			folder_name_arr := strings.Split(old_url, "/")
			folder_name = folder_name_arr[len(folder_name_arr)-1]
			if strings.HasSuffix(folder_name, ".git") {
				folder_name = strings.Split(folder_name, ".git")[0]
			}
			fmt.Println("Folder name:", folder_name)
		}
	}

	args[0] = getGitFile()
	err := RunCommand(args[0], args[1:]...)
	if err != nil || len(new_url) == 0 || len(folder_name) == 0 {
		panic(err)
	}
	fmt.Println("Clone done!!!")

	err = os.Chdir(folder_name)
	if err != nil {
		panic(err)
	}

	restoreCmd := "remote set-url origin " + old_url
	err = RunCommand(args[0], strings.Fields(restoreCmd)...)
	if err != nil {
		panic(err)
	}

	fmt.Println("Set remote done!!!")
}

func main() {
	cmdArgs := os.Args
	fmt.Println(cmdArgs)

	if len(cmdArgs) > 2 && cmdArgs[1] == "clone" {
		GgitClone(cmdArgs)
	} else {
		log.Fatal("请输入ggit clone https://github.com/xxx/xxx.git这样的命令格式")
	}
}
