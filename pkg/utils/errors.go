package utils

import (
	"fmt"
	"os"
	"strings"
)

func CheckErr(err error) {
	if err == nil {
		return
	}
	msg := err.Error()
	if !strings.HasPrefix(msg, "error: ") {
		msg = fmt.Sprintf("error: %s", msg)
	}
	fmt.Fprint(os.Stderr, msg)
	os.Exit(1)
}
