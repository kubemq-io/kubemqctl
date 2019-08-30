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
	"os/signal"
	"sync"
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
	ctx, cancel := context.WithCancel(ctx)
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting wd: %s", err)
	}
	opts := []getter.ClientOption{}
	opts = append(opts, getter.WithProgress(defaultProgressBar))
	client := &getter.Client{
		Ctx:     ctx,
		Src:     "https://github.com/kubemq-io/kubetools/releases/download/latest/kubetools_darwin_amd64",
		Dst:     "./dist",
		Pwd:     pwd,
		Mode:    getter.ClientModeAny,
		Options: opts,
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	errChan := make(chan error, 2)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := client.Get(); err != nil {
			errChan <- err
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	select {
	case sig := <-c:
		signal.Reset(os.Interrupt)
		cancel()
		wg.Wait()
		log.Printf("signal %v", sig)
	case <-ctx.Done():
		wg.Wait()
		log.Printf("success!")
	case err := <-errChan:
		wg.Wait()
		log.Fatalf("Error downloading: %s", err)
	}
	return nil
}
