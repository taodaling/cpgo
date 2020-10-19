package main

import "fmt"

func Red(content string) string {
	return fmt.Sprintf("\033[41;37m%v\033[0m", content)
}

func Green(content string) string {
	return fmt.Sprintf("\033[42;37m%v\033[0m", content)
}

func Purple(content string) string {
	return fmt.Sprintf("\033[45;37m%v\033[0m", content)
}
