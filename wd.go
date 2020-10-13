package main

import (
	"io/ioutil"
	"strings"
)

func Suffix(a string) string {
	index := strings.LastIndex(a, ".")
	if index == -1 {
		return ""
	}
	return a[index+1:]
}

func CollectFile(dirName string) (map[string][]string, error) {
	ans := make(map[string][]string)
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		ans[Suffix(filename)] = append(ans[Suffix(filename)], filename)
	}
	return ans, nil
}
