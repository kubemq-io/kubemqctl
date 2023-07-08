package main

import (
	"embed"
	"github.com/kubemq-io/kubemqctl/cmd/root"
	"github.com/kubemq-io/kubemqctl/web"
	"os"
)

//go:embed assets/*
var staticAssets embed.FS

var version string

func main() {
	web.StaticAssets = staticAssets
	root.Execute(version, os.Args)
}
