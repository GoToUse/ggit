package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func checkRepo(args []string) {
	gitRepo := args[0]
	if strings.HasSuffix(gitRepo, DefaultGithubSuffix) {
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
	} else {
		log.Fatalf("[Error]: github repo url doesn't end with `%s` suffix.", DefaultGithubSuffix)
	}
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
	fmt.Println("[folderAbsPath]", folderAbsPath)
	folderAbsPathStats, err := os.Stat(folderAbsPath)
	if err != nil {
		return false
	}
	if folderAbsPathStats.IsDir() {
		return true
	}
	return false
}
