package utils

import (
	"fmt"
	"github.com/gookit/color"
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
	fmt.Print(Title(msg))
}
func Printf(format string, args ...interface{}) {
	fmt.Print(Title(fmt.Sprintf(format, args...)))
}
func Printlnf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	color.Print(Title(fmt.Sprintf("%s\n", str)))
}
func PrintlnfNoTitle(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	color.Print(fmt.Sprintf("%s\n", str))
}
func PrintAndExit(msg string) {
	fmt.Println(Title(msg))
}

func PrintfAndExit(format string, args ...interface{}) {
	fmt.Println(Title(fmt.Sprintf(format, args...)))
}
