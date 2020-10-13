package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func fetchWords(s string) []string {
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words
}

func checkResult(expect, actual string) error {
	expectWords := fetchWords(expect)
	actualWords := fetchWords(actual)
	if len(expectWords) != len(actualWords) {
		return errors.New("DIFFERENT LENGTH")
	}
	for i := 0; i < len(expectWords); i++ {
		if expectWords[i] != actualWords[i] {
			return errors.New(strconv.Itoa(i+1) + "TH ELEMENT IS DIFFERENT")
		}
	}
	return nil
}

func run(testName, progName, input, output string, ch chan<- string) {
	res := new(bytes.Buffer)

	fmt.Fprintf(res, "TESTCASE %v:", testName)
	outbuf := new(bytes.Buffer)
	errbuf := new(bytes.Buffer)
	cmd := exec.Command(progName)
	cmd.Stdin = strings.NewReader(input)
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf
	time.AfterFunc(10*time.Second, func() {
		syscall.Kill(cmd.Process.Pid, syscall.SIGKILL)
	})
	now := time.Now()
	err := cmd.Run()
	delta := time.Now().Sub(now).Milliseconds()
	fmt.Fprintf(res, "%vms:", delta)
	if err != nil {
		fmt.Fprintf(res, "ERROR:%v", delta, err.Error())
	} else {
		err = checkResult(output, outbuf.String())
		if err != nil {
			fmt.Fprintf(res, "WA:%v", err.Error())
		} else {
			fmt.Fprintf(res, "PASS")
		}
	}
	ch <- res.String()
}

func runEntry(progName, inputFile, outputFile string, ch chan<- string) {
	input, _ := ioutil.ReadFile(inputFile)
	output, _ := ioutil.ReadFile(outputFile)
	run(inputFile, progName, string(input), string(output), ch)
}

func TestEntry(wd string) {
	lastModify := time.Now()
	progName := wd + string(os.PathSeparator) + "main"
	for {
		time.Sleep(time.Second)

		if stat, err := os.Stat(progName); err != nil || !stat.ModTime().After(lastModify) {
			continue
		} else {
			lastModify = stat.ModTime()
		}

		if filesBySuffix, err := CollectFile(wd); err != nil {
			log.Println("\nCan't collect file:", err)
			continue
		} else {
			fmt.Fprintln(os.Stdout, "Detect program has changed, run tests...")
			inputs := filesBySuffix["in"]
			outputs := filesBySuffix["out"]
			sort.Strings(inputs)
			sort.Strings(outputs)
			if len(inputs) != len(outputs) {
				log.Println("The number of input and output doesn't match")
				continue
			}
			var chans []chan string
			for i := 0; i < len(inputs); i++ {
				ch := make(chan string)
				chans = append(chans, ch)
				go runEntry(progName, inputs[i], outputs[i], ch)
			}
			for _, ch := range chans {
				msg := <-ch
				fmt.Fprintln(os.Stdout, msg)
			}
			fmt.Fprintln(os.Stdout, "All test has done!\n")
		}
	}
}