package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func loadFile(added map[string]bool, path string, includePaths []string) string {
	log.Println("load file", path)
	path, _ = filepath.Abs(path)
	if added[path] {
		return ""
	}
	added[path] = true
	headerContent, _ := ioutil.ReadFile(path)

	//replace some unused token
	pragma, _ := regexp.Compile(`#pragma once`)
	replaced := pragma.ReplaceAllString(string(headerContent), "")

	return rec(replaced, added, filepath.Dir(path), includePaths)
}

func existFile(path string)bool{
		s, err := os.Stat(path)
		return err == nil && !s.IsDir()
}

func rec(content string, added map[string]bool, wd string, includePaths []string) string {
	includeRe, _ := regexp.Compile(`#include\s*"(.*)"`)
	return includeRe.ReplaceAllStringFunc(content, func(match string) string {
		sub := includeRe.FindSubmatch([]byte(match))
		headerName := string(sub[1])
		//don't inline those files
		if headerName == "prettyprint.h" {
			return ""
		}
		if existFile(wd+string(os.PathSeparator)+headerName) {
			return loadFile(added, wd+string(os.PathSeparator)+headerName, includePaths)
		}
		for _, includePath := range includePaths{
			if existFile(includePath+headerName) {
				return loadFile(added, includePath+headerName, includePaths)
			}
		}
		fmt.Fprintln(os.Stdout, "FILE " + headerName + " NOT EXIST!")
		return `#include "` + headerName + `"`
	})
}

func Inline(wd string) {
	cppIncludePathEnv := os.Getenv("CPLUS_INCLUDE_PATH")
	includePaths := make([]string, 0)
	includePaths = append(includePaths, wd + string(os.PathSeparator))
	includePaths = append(includePaths, "")
	var cppIncludePaths []string
	sysType := runtime.GOOS
	if sysType == "windows" {
        cppIncludePaths = strings.Split(cppIncludePathEnv, ";")
    }else{
		cppIncludePaths = strings.Split(cppIncludePathEnv, ":")
	}
	for _, path := range cppIncludePaths {
		if path != "" {
			includePaths = append(includePaths, path + string(os.PathSeparator))
		}
	}

	merge := loadFile(make(map[string]bool), wd+string(os.PathSeparator)+"main.cpp", includePaths)
	ioutil.WriteFile(wd+string(os.PathSeparator)+"inline.cpp", []byte(merge), os.FileMode(0777))
	fmt.Fprintln(os.Stdout, "INLINE SUCCESS, PLEASE VIEW inline.cpp")
}
