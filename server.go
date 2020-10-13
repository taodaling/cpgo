package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Test struct {
	Input, Output string
}

type Body struct {
	Name, Group, Url, TestType string
	MemoryLimit, TimeLimit     int
	Interactive                bool
	Tests                      []Test
}

func ServerEntry(wd string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Receive request", r)
		var body Body
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Read post data error:", err)
			return
		}
		log.Println("Read post content:", string(content))
		err = json.Unmarshal(content, &body)
		if err != nil {
			log.Println("Can't parse the body:", err)
			return
		}
		log.Println("json result:", body)

		if filesBySuffix, err := CollectFile(wd); err != nil {
			log.Println("Can't collect file:", err)
			return
		} else {
			candidates := make([]string, 0)
			candidates = append(candidates, filesBySuffix["in"]...)
			candidates = append(candidates, filesBySuffix["out"]...)
			for _, item := range candidates {
				err = os.Remove(item)
				if err != nil {
					log.Println("Can't remove file:", err)
					return
				}
			}
		}

		//create test
		for i, test := range body.Tests {
			ioutil.WriteFile(strconv.Itoa(i)+".in", []byte(test.Input), os.FileMode(0777))
			ioutil.WriteFile(strconv.Itoa(i)+".out", []byte(test.Output), os.FileMode(0777))
		}
	})

	log.Println("Server booting...")
	log.Fatal(http.ListenAndServe(":50823", nil))
}
