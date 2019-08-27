package utils

import (
	"fmt"
)

func Println(msg string) {
	fmt.Println(msg)
}

func Print(msg string) {
	fmt.Print(msg)
}
func Printf(format string, args ...interface{}) {
	fmt.Print(fmt.Sprintf(format, args...))
}

func Printlnf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func PrintAndExit(msg string) {
	fmt.Println(msg)
}

func PrintfAndExit(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}
