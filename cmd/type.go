package cmd

import "time"

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

type (
	GitS struct {
		FilePath  string
		Website   string
		UrlSuffix string
	}
	MirrorUrlS map[string][]string
)

type CallBack func() error
