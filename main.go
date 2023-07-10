package main

import (
	"embed"
	"github.com/kubemq-io/kubemqctl/cmd/root"
	"github.com/kubemq-io/kubemqctl/web"
	"os"
)

//go:embed assets-web/*
var staticAssetsWeb embed.FS

//go:embed assets-builder/*
var staticAssetsBuilder embed.FS
var version string

func main() {
	web.StaticAssetsWeb = staticAssetsWeb
	web.StaticAssetsBuilder = staticAssetsBuilder
	root.Execute(version, os.Args)
}
