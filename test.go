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

func red(content string) string {
	return fmt.Sprintf("\033[31m%v\033[0m", content)
}

func green(content string) string {
	return fmt.Sprintf("\033[32m%v\033[0m", content)
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
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGKILL)
		}
	})
	now := time.Now()
	err := cmd.Run()
	delta := time.Now().Sub(now).Milliseconds()
	fmt.Fprintf(res, "%vms:", delta)
	if err != nil {
		fmt.Fprintf(res, "\033[35mERROR\033[0m:%v", delta, err.Error())
	} else {
		err = checkResult(output, outbuf.String())
		if err != nil {
			fmt.Fprintf(res, "\033[31mWA\033[0m:%v", err.Error())
		} else {
			fmt.Fprintf(res, "\033[32mPASS\033[0m")
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
			fmt.Fprint(os.Stdout, "\033[2J\033[0;0H")
			fmt.Fprintf(os.Stdout, "[%v]Detect program has changed, run tests...\n", time.Now().Format("15:04:05"))
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
			//inline file at last
			Inline(wd)
		}
	}
}
