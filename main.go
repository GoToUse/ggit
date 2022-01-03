package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-ping/ping"
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
	SortHost []rttHost
	Args     []string
	// GitRepoInfo git-repo struct
	GitRepoInfo struct {
		RawPath  string
		Author   string
		RepoName string
	}
	// CheckRepo interface
	CheckRepo interface {
		checkIsAGitRepo() bool
		folderDuplicate() bool
	}
)

var (
	wg                       sync.WaitGroup
	wg1                      sync.WaitGroup
	GitRepoInit              = new(GitRepoInfo)
	DEFAULT_MIRROR_URL_ARRAY = []string{
		"https://hub.fastgit.org/",
		"https://github.com.cnpmjs.org/",
		"https://gitclone.com/",
		"https://github.wuyanzheshui.workers.dev/",
	}
)

func init() {
	cmdArgs := os.Args
	cmdArgsLength := len(cmdArgs)
	if cmdArgsLength > 2 &&
		cmdArgs[1] == "clone" &&
		strings.HasSuffix(cmdArgs[2], DEFAULT_GITHUB_SUFFIX) {
		gitRepo := cmdArgs[2]
		GitRepoInit.RawPath = gitRepo
		argsTwoList := strings.Split(gitRepo, "/")
		argsTwoRepo := argsTwoList[len(argsTwoList)-1]
		argsTwoAuthor := argsTwoList[len(argsTwoList)-2]
		folderName := strings.Split(argsTwoRepo, ".git")[0]
		GitRepoInit.Author = argsTwoAuthor
		GitRepoInit.RepoName = folderName

		// Exit if the folder is already a git-repo.
		var checkRepoF CheckRepo
		checkRepoF = GitRepoInit
		if checkRepoF.checkIsAGitRepo() {
			log.Fatalf(
				"💨 %s already exists and is a git repository. Program will exit now...",
				GitRepoInit.RepoName,
			)
		}
		if checkRepoF.folderDuplicate() {
			log.Fatalf(
				"💨 Your current path already has a directory with the same name[%s]. Program will exit now...",
				GitRepoInit.RepoName,
			)
		}
	} else if cmdArgsLength > 3 {
		log.Fatal("[Error]:请输入[ggit clone https://github.com/xxx/xxx.git]这样的命令格式")
	} else {
		log.Fatal("[Error]:请输入[ggit clone https://github.com/xxx/xxx.git]这样的命令格式")
	}
}

func RunCommand(name string, args ...string) error {
	fmt.Println("Command:", append([]string{name}, args...))
	seperator := center(strings.ToUpper(args[0]), 40, "-")
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
	var oldUrl, newUrl string

	if strings.HasPrefix(args[2], DEFAULT_GITHUB_URL) {
		oldUrl = args[2]
		// 特别处理
		u, err := url.Parse(mirrorUrl)
		if err != nil {
			log.Panicf("%s is wrong, see details[%s]", mirrorUrl, err.Error())
		}
		if strings.Contains(mirrorUrl, "https://gitclone.com/") {
			// Check the git-repo if exists on the gitclone.com.
			if existOnGitClone(GitRepoInit.RepoName, GitRepoInit.Author) {
				ref, _ := u.Parse("github.com")
				githubCloneUrl := fmt.Sprintf("%s/", ref)
				newUrl = strings.ReplaceAll(oldUrl, DEFAULT_GITHUB_URL, githubCloneUrl)
			} else {
				newUrl = strings.ReplaceAll(oldUrl, DEFAULT_GITHUB_URL, mirrorUrl)
			}
		} else {
			newUrl = strings.ReplaceAll(oldUrl, DEFAULT_GITHUB_URL, mirrorUrl)
		}
		args[2] = newUrl
		fmt.Println("Folder name:", GitRepoInit.RepoName)
	} else {
		fmt.Printf("DEBUG: args[2]: %s\n", args[2])
		log.Fatal("github仓库地址有误, 请检查是否符合 [https://github.com/xxx/xxx.git] 标准路径.")
	}

	args[0] = getGitFile()
	err := RunCommand(args[0], args[1:]...)
	if err != nil || len(newUrl) == 0 || len(GitRepoInit.RepoName) == 0 {
		retryErr := Retry(3, 3*time.Second, func() error {
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

	cdr, _ := os.Getwd()
	repoAbsPath := path.Join(cdr, GitRepoInit.RepoName)
	err = os.Chdir(repoAbsPath)
	if err != nil {
		panic(err)
	}

	restoreCmd := "remote set-url origin " + oldUrl
	err = RunCommand(args[0], strings.Fields(restoreCmd)...)
	if err != nil {
		retryErr := Retry(3, 3*time.Second, func() error {
			fErr := RunCommand(args[0], args[1:]...)
			return fErr
		})
		if retryErr != nil {
			// 如果当前url不能正常工作，那么初始化args[2]的值
			args[2] = oldUrl
			// TODO: if error, delete this folder.
			err = os.RemoveAll(repoAbsPath)
			if err != nil {
				log.Fatalf("Remove wrong: %s", err.Error())
			}
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
	seperator := center("\U0001F973 Sort By Ping RTT Value \U0001F973", 80, "#")
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
	return rttMapList
}

func GgitClone(args Args) {
	var initTimes int
	sortHostRes := sortHost(DEFAULT_MIRROR_URL_ARRAY)
	fmt.Printf("Sorted list: %v\n%s\n", sortHostRes, strings.Repeat("*", 80))
	for _, v := range sortHostRes {
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

// center like `str.center` function in python.
func center(s string, n int, fill string) string {
	sLen := len(s)
	div := (n - sLen) / 2
	return strings.Repeat(fill, div) + fmt.Sprintf(" %s ", s) + strings.Repeat(fill, div)
}

func (g *GitRepoInfo) checkIsAGitRepo() bool {
	gitP := lookGitPath()
	rootP, _ := os.Getwd()
	destP := path.Join(rootP, g.RepoName)
	err := os.Chdir(destP)
	if err != nil {
		return false
	}
	err = RunCommand(gitP, []string{"rev-parse", "--is-inside-work-tree"}...)

	return err == nil
}

func (g *GitRepoInfo) folderDuplicate() bool {
	currentPwd, _ := os.Getwd()
	folderAbsPath := path.Join(currentPwd, g.RepoName)
	fmt.Println("folderAbsPath", folderAbsPath)
	folderAbsPathStats, err := os.Stat(folderAbsPath)
	if err != nil {
		return false
	}
	if folderAbsPathStats.IsDir() {
		return true
	}
	return false
}

func existOnGitClone(gitRepoName, gitAuthorName string) bool {
	queryPath := fmt.Sprintf("https://www.gitclone.com/gogs/search/clonesearch?q=%s", gitRepoName)
	res, err := http.Get(queryPath)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return false
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}
	// cannot search the repo
	noResult := doc.Find("p.clonesearch-noresult")
	if noResult.Length() != 0 {
		fmt.Println("Tips😅: gitclone.com didn't have this repo. Maybe you can add your git-repo to gitclone.com manually and then you can use it later. Details see website: https://gitclone.com/")
		return false
	}
	// find some thing
	repoExistList := make([]bool, 0, 1)
	doc.Find("div.item").Each(func(i int, s *goquery.Selection) {
		title := s.Find("a.name").Text()
		// pass the situation of empty strings
		if len(title) != 0 {
			// trimTitle's format: authorName/repoName
			trimTitle := strings.TrimSpace(title)
			// Concat the current git-repo format to compare with trimTitle
			currentGitInfo := gitAuthorName + "/" + gitRepoName
			if currentGitInfo == trimTitle {
				// true, it indicates the repo has been already backup in the gitclone.com.
				repoExistList = append(repoExistList, true)
				return
			}
		}
	})

	// Check the git-repo if exists
	if len(repoExistList) != 0 && repoExistList[0] {
		return true
	} else {
		return false
	}
}

func main() {
	cmdArgs := os.Args
	fmt.Println(cmdArgs)
	GgitClone(cmdArgs)
}
