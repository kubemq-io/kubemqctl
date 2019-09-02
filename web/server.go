package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-getter"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/skratchdot/open-golang/open"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ServerOptions struct {
	Cfg  *config.Config
	Port int
	Path string
}

type Config struct {
	DashboardApi string `json:"dashboard_api"`
	SocketApi    string `json:"socket_api"`
}

func createConfigAndSave(path string, apiPort int, restPort int) error {
	c := &Config{
		DashboardApi: fmt.Sprintf("http://localhost:%d/v1/stats/", apiPort),
		SocketApi:    fmt.Sprintf("ws://localhost:%d/v1/stats/", apiPort),
	}
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (s *ServerOptions) Run(ctx context.Context) error {
	err := createConfigAndSave(s.Path+"/config.json", s.Cfg.ApiPort, s.Cfg.RestPort)
	if err != nil {
		return err
	}
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
