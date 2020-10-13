package main

import (
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stderr)
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Getwd error:", err)
	}
	go ServerEntry(wd)
	TestEntry(wd)
}
