package main

import (
	"github.com/kubemq-io/kubemqctl/cmd/root"
	"os"
)

var version string

func main() {
	root.Execute(version, os.Args)
}
