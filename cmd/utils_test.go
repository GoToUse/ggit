package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindKey(t *testing.T) {
	testCases1 := map[string][]string{
		"none": {
			"https://hub.fastgit.xyz/",
			"https://github.com.cnpmjs.org/",
			"https://hub.連接.台灣/",
		},
		"https://github.com": {
			"https://ghproxy.com/",
			"https://www.github.do/",
		},
	}

	equalMsg := "they should be equal."

	got1 := FindKey(testCases1, []string{
		"https://ghproxy.com/",
		"https://www.github.do/",
	})
	want1 := "https://github.com"
	assert.Equal(t, want1, got1, equalMsg)

	got2 := FindKey(testCases1, []string{
		"https://hub.fastgit.xyz/",
		"https://github.com.cnpmjs.org/",
		"https://hub.連接.台灣/",
	})
	want2 := "none"
	assert.Equal(t, want2, got2, equalMsg)

	got3 := FindKey(testCases1, "https://www.github.do/")
	want3 := "https://github.com"
	assert.Equal(t, want3, got3, equalMsg)

	got4 := FindKey(testCases1, "https://hub.連接.台灣/")
	want4 := "none"
	assert.Equal(t, want4, got4, equalMsg)

	got5 := FindKey(testCases1, "https://gitclone.com/")
	want5 := "" // 不存在的返回空字符串
	assert.Equal(t, want5, got5, equalMsg)
}
