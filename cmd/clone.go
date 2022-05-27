package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone <git repo url>",
	Short: "Clone the specified repo from github.com",
	Run: func(cmd *cobra.Command, args []string) {
		other, err := cmd.Flags().GetStringArray("other")
		if err != nil {
			os.Exit(1)
		}
		args = append(args, other...)
		checkRepo(args)
		fullParams := func(arr []string) []string {
			for i, v := range map[int]string{0: os.Args[0], 1: os.Args[1]} {
				arr = insertElem(arr, i, v)
			}
			return arr
		}(args)
		fmt.Println("fullParams:", fullParams)
		GgitClone(fullParams)
	},
}

func init() {
	cloneCmd.Flags().StringArrayP(
		"other",
		"o",
		[]string{},
		"other sub-commands of clone-command in git. \nWrap it in double quotation marks. \neg. \"--depth=1\"",
	)
}

var GitRepoInit = new(GitRepoInfo)

func ggitClone(args Args, mirrorUrl string) error {
	var oldUrl, newUrl, ref, githubCloneUrl string

	if strings.HasPrefix(args[2], DefaultGithubUrl) {
		oldUrl = args[2]
		// 特别处理
		u, err := url.Parse(mirrorUrl)

		if err != nil {
			log.Panicf("%s is wrong, see details[%s]", mirrorUrl, err.Error())
		}

		concatString := FindKey(DefaultMirrorUrlMap, mirrorUrl)
		ref = strings.Join([]string{strings.TrimSuffix(u.String(), "/"), concatString}, "/")
		githubCloneUrl = fmt.Sprintf("%s/", ref)

		if concatString == "" {
			newUrl = oldUrl
		} else if concatString == "none" {
			newUrl = strings.ReplaceAll(oldUrl, DefaultGithubUrl, mirrorUrl)
		} else {
			newUrl = strings.ReplaceAll(oldUrl, DefaultGithubUrl, githubCloneUrl)

			if concatString == "gitclone.com" && !existOnGitClone(GitRepoInit.RepoName, GitRepoInit.Author) {
				newUrl = oldUrl
			}
		}

		args[2] = newUrl
		fmt.Println("Folder name:", GitRepoInit.RepoName)
	} else {
		fmt.Printf("DEBUG: args[2]: %s\n", args[2])
		log.Fatal("github仓库地址有误, 请检查是否符合 [https://github.com/xxx/xxx.git] 标准路径.")
	}

	return nil

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
		log.Fatalf("os.Chdir err: %v", err)
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

func GgitClone(args Args) {
	var initTimes int
	sortHostRes := sortHost(HostValues(DefaultMirrorUrlMap))

	fmt.Println(center("Sorted list", 80, "*"))
	RenderTable(sortHostRes)
	fmt.Println(strings.Repeat("*", 80))

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

	if initTimes == len(DefaultMirrorUrlMap) {
		log.Fatal("Sorry: All mirrors are unusable.")
	}
}

func existOnGitClone(gitRepoName, gitAuthorName string) bool {
	queryPath := fmt.Sprintf("https://www.gitclone.com/gogs/search/clonesearch?q=%s", gitRepoName)
	res, err := http.Get(queryPath)
	if err != nil {
		log.Fatalf("http.Get err: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return false
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("goquery.NewDocumentFromReader err: %v", err)
	}
	// cannot search the repo
	noResult := doc.Find("p.clonesearch-noresult")
	if noResult.Length() != 0 {
		fmt.Println("Tips😅: gitclone.com didn't have this repo. Maybe you can add your git-repo to gitclone.com manually and then you can use it later. Details see website: https://gitclone.com/")
		return false
	}
	// find something
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
