package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"runtime"
)

const SUCCESS = 0
const WA = 1
const ERROR = 2

type TestResult struct {
	Input, Output, Actual, Msg string
	State                      int
}

func (t *TestResult) Format(w io.Writer) {
	var content = t.Msg
	if t.State == SUCCESS {
		content = Green(content)
	} else if t.State == WA {
		content = Red(content)
	} else if t.State == ERROR {
		content = Purple(content)
	}
	fmt.Fprintln(w, content)
}

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

func run(testName, progName, input, output string, ch chan<- TestResult) {
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
	actual := outbuf.String()
	result := TestResult{
		Input:  input,
		Output: output,
		Actual: actual,
	}
	fmt.Fprintf(res, "%vms:", delta)
	if err != nil {
		fmt.Fprintf(res, "ERROR:%v", delta, err.Error())
		result.State = ERROR
	} else {
		err = checkResult(output, actual)
		if err != nil {
			fmt.Fprintf(res, "WA:%v", err.Error())
			result.State = WA
		} else {
			fmt.Fprintf(res, "PASS")
			result.State = SUCCESS
		}
	}

	result.Msg = res.String()
	ch <- result
}

func RunEntry(progName, inputFile, outputFile string, ch chan<- TestResult) {
	input, _ := ioutil.ReadFile(inputFile)
	output, _ := ioutil.ReadFile(outputFile)
	run(inputFile, progName, string(input), string(output), ch)
}

func TestEntry(wd string) {
	lastModify := time.Now()
	progName := wd + string(os.PathSeparator) + "main"
	sysType := runtime.GOOS
	if sysType == "windows" {
        progName += ".exe"
    }

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
				fmt.Fprintf(os.Stdout, "The number of input and output doesn't match")
				continue
			}
			var chans []chan TestResult
			for i := 0; i < len(inputs); i++ {
				ch := make(chan TestResult)
				chans = append(chans, ch)
				go RunEntry(progName, inputs[i], outputs[i], ch)
			}
			for _, ch := range chans {
				result := <-ch
				result.Format(os.Stdout)
			}
			fmt.Fprintln(os.Stdout, "All test has done!\n")
			//inline file at last
			Inline(wd)
		}
	}
}
