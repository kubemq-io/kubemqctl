package main

import (
	"github.com/kubemq-io/kubetools/cmd/root"
)

var version string

func main() {
	root.Execute(version)
}
