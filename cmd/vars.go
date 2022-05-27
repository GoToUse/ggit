package cmd

import "sync"

var (
	DefaultGitPath      string
	DefaultGithubUrl    string
	DefaultGithubSuffix string
	DefaultMirrorUrlMap map[string][]string
)

var (
	wg  sync.WaitGroup
	wg1 sync.WaitGroup
)

var (
	GitC         GitS
	mirrorUrlArr MirrorUrlS
)
