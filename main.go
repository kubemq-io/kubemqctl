package main

import (
	"github.com/kubemq-io/kubemqctl/cmd/root"
)

var version = "v2.0.0rc1"

func main() {
	root.Execute(version)
}
