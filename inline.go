package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func loadFile(added map[string]bool, path string) string {
	log.Println("load file", path)
	path, _ = filepath.Abs(path)
	if added[path] {
		return ""
	}
	added[path] = true
	headerContent, _ := ioutil.ReadFile(path)
	return rec(string(headerContent), added, filepath.Dir(path))
}

func rec(content string, added map[string]bool, wd string) string {
	includeRe, _ := regexp.Compile(`#include\s*"(.*)"`)
	return includeRe.ReplaceAllStringFunc(content, func(match string) string {
		sub := includeRe.FindSubmatch([]byte(match))
		headerName := string(sub[1])
		return loadFile(added, wd+string(os.PathSeparator)+headerName)
	})
}

func Inline(wd string) {
	log.Println("Now inline work")
	merge := loadFile(make(map[string]bool), wd+string(os.PathSeparator)+"main.cpp")
	ioutil.WriteFile(wd+string(os.PathSeparator)+"inline.cpp", []byte(merge), os.FileMode(0777))
	fmt.Fprintln(os.Stdout, "INLINE SUCCESS, PLEASE VIEW inline.cpp")
}
