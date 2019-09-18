package dashboard

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/kubemq-io/kubemqctl/web"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type DashboardOptions struct {
	cfg    *config.Config
	update bool
}

var dashboardExamples = `
	# Run KubeMQ dashboard web view
	kubemqctl dashboard

	# Run KubeMQ dashboard and update version
	kubemqctl dashboard -u
`
var dashboardLong = `Dashboard command allows to start a web view of KubeMQ cluster dashboard`
var dashboardShort = `Run KubeMQ dashboard web view command`

// NewCmdCreate returns new initialized instance of create sub query
func NewCmdDashboard(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := DashboardOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"web", "dash"},
		Short:   dashboardLong,
		Long:    dashboardShort,
		Example: dashboardExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().BoolVarP(&o.update, "update", "u", false, "update dashboard version")
	return cmd
}

func (o *DashboardOptions) Complete(args []string) error {

	return nil
}

func (o *DashboardOptions) Validate() error {
	return nil
}
func exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func (o *DashboardOptions) Run(ctx context.Context) error {
	s := &web.ServerOptions{
		Cfg:  o.cfg,
		Port: 6700,
		Path: "./dashboard",
	}
	if !o.update {
		if !exists("./dashboard/index.html") {
			o.update = true
		}
	}
	if o.update {
		err := s.Download(ctx)
		if err != nil {
			utils.CheckErr(err)
		}
	}
	err := o.setConnections(s.Path)
	if err != nil {
		return err
	}
	err = s.Run(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (o *DashboardOptions) setConnections(path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "main") {
			data, err := ioutil.ReadFile(filepath.Join(path, f.Name()))
			if err != nil {
				return err
			}
			file := string(data)
			file = strings.Replace(file, "DASHBOARD_API_PLACEMENT", o.cfg.GetApiHttpURI()+"/v1/stats", -1)
			file = strings.Replace(file, "SOCKET_API_PLACEMENT", o.cfg.GetApiWsURI()+"/v1/stats", -1)
			err = ioutil.WriteFile(filepath.Join(path, f.Name()), []byte(file), 0644)
			if err != nil {
				return err
			}

			return nil

		}
	}
	return fmt.Errorf("invalid dashbord distribution content")
}
