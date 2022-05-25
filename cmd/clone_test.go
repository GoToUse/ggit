package cmd

import (
	"strings"
	"testing"
)

func TestRetrieveHost(t *testing.T) {
	urlList := []string{
		"https://hub.連接.台灣/",
		"https://www.github.do/",
		"https://hub.おうか.tw/",
		"https://hub.fastgit.xyz/",
	}

	for _, url := range urlList {
		host := retrieveHost(url)
		wanted := strings.Trim(strings.Split(url, "//")[1], "/")

		if host != wanted {
			t.Fatalf("retrieveHost(%s) = %s; wanted (%s)", url, host, wanted)
		}
	}
}
