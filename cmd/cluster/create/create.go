package create

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"time"

	conf "github.com/kubemq-io/kubemqctl/pkg/k8s/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"

	apiv1 "k8s.io/api/core/v1"
	"os"
	"strings"
)

type CreateOptions struct {
	cfg           *config.Config
	setOptions    bool
	exportFile    bool
	file          string
	importData    string
	optionsMenu   *conf.Menu
	deployOptions *deployment.Options
}

var createExamples = `
	# Create default KubeMQ cluster
	# kubemqctl cluster create -t b33600cc-93ef-4395-bba3-13131eb27d5e 

	# Import KubeMQ cluster yaml file  
	# kubemqctl cluster create -f kubemq-cluster.yaml

	# Create KubeMQ cluster with options
	# kubemqctl cluster create -t b3330scc-93ef-4395-bba3-13131sb2785e -o

	# Export KubeMQ cluster yaml file    
	# kubemqctl cluster create -t b3330scc-93ef-4395-bba3-13131sb2785e -e 
`
var createLong = `Create a KubeMQ cluster`
var createShort = `Create a KubeMQ cluster`

func NewCmdCreate(cfg *config.Config) *cobra.Command {
	o := &CreateOptions{
		cfg:           cfg,
		optionsMenu:   conf.CreateMenu,
		deployOptions: &deployment.Options{},
	}
	cmd := &cobra.Command{

		Use:     "create",
		Aliases: []string{"c", "cr"},
		Short:   createShort,
		Long:    createLong,
		Example: createExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().StringVarP(&o.deployOptions.Token, "token", "t", "", "set KubeMQ Token")
	cmd.PersistentFlags().BoolVarP(&o.setOptions, "options", "o", false, "create KubeMQ cluster with options")
	cmd.PersistentFlags().BoolVarP(&o.exportFile, "export", "e", false, "generate yaml configuration file")
	cmd.PersistentFlags().StringVarP(&o.file, "file", "f", "", "import configuration yaml file")

	return cmd
}

func (o *CreateOptions) Complete(args []string) error {

	if o.file != "" {
		buff, err := ioutil.ReadFile(o.file)
		if err != nil {
			return err
		}
		o.importData = string(buff)
		return nil
	}

	if o.deployOptions.Token == "" {
		if o.cfg.DefaultToken != "" {
			o.deployOptions.Token = o.cfg.DefaultToken
			utils.Printlnf("Create using save KubeMQ token: %s", o.cfg.DefaultToken)
		} else {
			return fmt.Errorf("No KubeMQ token provided")
		}

	} else {
		o.cfg.DefaultToken = o.deployOptions.Token
		_ = o.cfg.Save()
	}

	if o.setOptions {
		err := conf.CreateMenu.Run()
		if err != nil {
			return err
		}
		o.deployOptions.AppsVersion = "apps/v1"
		o.deployOptions.CoreVersion = "v1"
		o.deployOptions.Name = conf.CreateBasicOptions.Name
		o.deployOptions.Namespace = conf.CreateBasicOptions.Namespace
		o.deployOptions.Version = conf.CreateBasicOptions.Image
		o.deployOptions.Replicas = conf.CreateBasicOptions.Replicas
		o.deployOptions.Volume = conf.CreateBasicOptions.Vol
		switch conf.CreateBasicOptions.ServiceMode {
		case "NodePort":
			o.deployOptions.IsNodePort = true
		case "LoadBalancer":
			o.deployOptions.IsLoadBalance = true
		}
		return nil
	}
	return o.setDefaultOptions()
}

func (o *CreateOptions) Validate() error {
	return nil
}

func (o *CreateOptions) Run(ctx context.Context) error {
	sd, err := deployment.NewStatefulSetDeployment(o.cfg)
	if err != nil {
		return err
	}

	if o.importData != "" {
		err := sd.Import(o.importData)
		if err != nil {
			return err
		}
	} else {
		var err error
		err = sd.CreateStatefulSetDeployment(o.deployOptions, o.optionsMenu)
		if err != nil {
			return err
		}
	}

	utils.Printlnf("Create started...")

	stsName := sd.StatefulSet.Name
	stsNamespace := sd.StatefulSet.Namespace
	if o.exportFile {

		f, err := os.Create(fmt.Sprintf("%s.yaml", stsName))
		if err != nil {
			return err
		}
		err = sd.Export(f)
		if err != nil {
			utils.Printlnf("export to file %s failed", stsName)
			return err
		}
		utils.Printlnf("export to file %s.yaml completed", stsName)
		return nil
	}

	executed, err := sd.Execute(sd.StatefulSet.Name, sd.StatefulSet.Namespace)
	if err != nil {
		return err
	}
	if !executed {
		return nil
	}
	utils.Printlnf("Create StatefulSet %s/%s progress:", stsNamespace, stsName)

	stsDone := make(chan struct{})
	stsChan := make(chan *appsv1.StatefulSet)

	evtDone := make(chan struct{})
	evtChan := make(chan *apiv1.Event)
	timeNow := time.Now()
	go sd.Client.GetStatefulSetEvents(ctx, stsChan, stsDone)
	go sd.Client.GetEvents(ctx, evtChan, evtDone)

	//w := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', tabwriter.TabIndent)
	//fmt.Fprintf(w, "Type\tReason\tMessage\n")
	//w.Flush()
	//w = tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', tabwriter.TabIndent)
	//	time.Sleep(2 * time.Second)

	for {
		select {
		case sts := <-stsChan:

			if sts.Name == sd.StatefulSet.Name && sts.Namespace == sd.StatefulSet.Namespace {

				rep := *sd.StatefulSet.Spec.Replicas
				utils.Printlnf("[StatefulSet] [Desired - %d] [Current - %d] [Ready - %d]", rep, sts.Status.Replicas, sts.Status.ReadyReplicas)
				//if rep == sts.Status.Replicas && sts.Status.Replicas == sts.Status.ReadyReplicas {
				//	//fmt.Fprintf(w, "%d\t%d\t%d\n", rep, sts.Status.Replicas, sts.Status.ReadyReplicas)
				//	//w.Flush()
				//	//utils.Printlnf("StatefulSet: Desired-%d Current-%d Ready-%d", rep, sts.Status.Replicas, sts.Status.ReadyReplicas)
				//	stsDone <- struct{}{}
				//	evtDone <- struct{}{}
				//	return nil
				//} else {
				//	//fmt.Fprintf(w, "%d\t%d\t%d\n", rep, sts.Status.Replicas, sts.Status.ReadyReplicas)
				//	//w.Flush()
				//	//utils.Printlnf("StatefulSet: Desired-%d Current-%d Ready-%d", rep, sts.Status.Replicas, sts.Status.ReadyReplicas)
				//}
			}
		case e := <-evtChan:
			if !strings.Contains(e.InvolvedObject.Name, stsName) {
				continue
			}
			if e.LastTimestamp.Sub(timeNow) < 0 {
				continue
			}
			if e.Count == 1 && time.Now().Sub(e.FirstTimestamp.Time) < time.Second {
				utils.Printlnf("%d [Event] [%s] [%s] [%s] [%s] -> %s", e.Count, e.Type, e.Reason, utils.TranslateTimestampSince(e.FirstTimestamp), e.InvolvedObject.Name, strings.TrimSpace(e.Message))

			}
			//var interval string
			//if e.Count > 1 {
			//	interval = fmt.Sprintf("%s (x%d over %s)", utils.TranslateTimestampSince(e.LastTimestamp), e.Count, utils.TranslateTimestampSince(e.FirstTimestamp))
			//} else {
			//	interval = utils.TranslateTimestampSince(e.FirstTimestamp)
			//}
			//utils.Printlnf("%d [Event] [%s] [%s] [%s] [%s] -> %s", e.Count, e.Type, e.Reason, interval, e.InvolvedObject.Name, strings.TrimSpace(e.Message))
		case <-ctx.Done():
			return nil
		}
	}

	//
	//<-ctx.Done()
	//return nil
}
func (o *CreateOptions) setDefaultOptions() error {

	o.deployOptions.AppsVersion = "apps/v1"
	o.deployOptions.CoreVersion = "v1"
	o.deployOptions.Name = "kubemq-cluster"
	o.deployOptions.Namespace = "default"
	o.deployOptions.Version = "latest"
	o.deployOptions.Replicas = 3
	o.deployOptions.Volume = 0
	utils.Printlnf("Create KubeMQ cluster with default options:")
	utils.Printlnf("\tKubeMQ Token: %s", o.deployOptions.Token)
	utils.Printlnf("\tCluster Name: %s", o.deployOptions.Name)
	utils.Printlnf("\tCluster Namespace: %s", o.deployOptions.Namespace)
	utils.Printlnf("\tCluster Docker Image: kubemq/kubemq:%s", o.deployOptions.Version)
	utils.Printlnf("\tCluster Replicas: %d", o.deployOptions.Replicas)
	utils.Printlnf("\tCluster PVC Size: %d", o.deployOptions.Volume)

	return nil
}
