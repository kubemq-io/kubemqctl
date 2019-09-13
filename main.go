package main

import (
	"github.com/kubemq-io/kubemqctl/cmd/root"
)

var version string

func main() {
	root.Execute(version)
}
