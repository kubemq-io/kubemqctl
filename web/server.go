package web

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-getter"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/skratchdot/open-golang/open"
	"log"
	"net/http"
	"os"
)

type ServerOptions struct {
	Cfg  *config.Config
	Port int
	Path string
}

func (s *ServerOptions) Run(ctx context.Context) error {

	fs := http.FileServer(http.Dir(s.Path))
	http.Handle("/", fs)

	go func() {

		path := fmt.Sprintf("http://localhost:%d/", s.Port)
		utils.Printlnf("start dashboard on '%s' ...", path)
		open.Run(path)
		err := http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil)
		if err != nil {
			utils.CheckErr(fmt.Errorf("dashboard runs already"))
		}

	}()

	return nil
}

func (s *ServerOptions) Download(ctx context.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting wd: %s", err)
	}
	opts := []getter.ClientOption{}
	opts = append(opts, getter.WithProgress(defaultProgressBar))
	client := &getter.Client{
		Ctx:              ctx,
		Src:              "https://github.com/kubemq-io/kubetools/releases/download/latest/dashboard.zip",
		Dst:              "./",
		Pwd:              pwd,
		Mode:             getter.ClientModeDir,
		Detectors:        nil,
		Decompressors:    nil,
		Getters:          nil,
		Dir:              false,
		ProgressListener: nil,
		Options:          opts,
	}

	utils.Println("download dashboard latest version...")
	err = client.Get()
	if err != nil {
		return err
	}

	utils.Println("downloaded...")

	return nil
}
