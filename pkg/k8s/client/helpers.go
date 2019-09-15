package client

import "strings"

func StringSplit(input string) (string, string) {
	pair := strings.Split(input, "/")
	if len(pair) == 2 {
		return pair[0], pair[1]
	}
	return "", ""
}
