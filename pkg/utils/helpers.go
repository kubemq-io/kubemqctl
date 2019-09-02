package utils

import (
	"fmt"
	"strings"
)

func Title(input string) string {
	words := strings.Fields(input)
	if len(words) > 0 {
		words[0] = strings.Title(words[0])
	}
	return strings.Join(words, " ")
}
func Println(msg string) {

	fmt.Println(Title(msg))
}

func Print(msg string) {
	fmt.Print(strings.Title(msg))
}
func Printf(format string, args ...interface{}) {
	fmt.Print(strings.Title(fmt.Sprintf(format, args...)))
}

func Printlnf(format string, args ...interface{}) {
	fmt.Println(Title(fmt.Sprintf(format, args...)))
}

func PrintAndExit(msg string) {
	fmt.Println(strings.Title(msg))
}

func PrintfAndExit(format string, args ...interface{}) {
	fmt.Println(strings.Title(fmt.Sprintf(format, args...)))
}
