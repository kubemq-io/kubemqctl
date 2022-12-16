package install

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"

	"github.com/spf13/cobra"
)

type InstallOptions struct {
	cfg         *config.Config
	name        string
	namespace   string
	key         string
	licenseFile string
	replicas    int
	licenseData string
	kubeCtlPath string
}

var installExamples = `
	# Execute install Kubemq cluster
	kubemqctl install -k <key> -n <namespace> -c <cluster name>	
`

var (
	installLong  = `Executes Kubemq install cluster command`
	installShort = `Executes Kubemq install cluster command`
)

func NewCmdInstall(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &InstallOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   installShort,
		Long:    installLong,
		Example: installExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().StringVarP(&o.name, "name", "c", "kubemq-cluster", "kubemq cluster name")
	cmd.PersistentFlags().StringVarP(&o.namespace, "namespace", "n", "kubemq", "kubemq cluster namespace")
	cmd.PersistentFlags().StringVarP(&o.key, "key", "k", "", "kubemq license key")
	cmd.PersistentFlags().StringVarP(&o.licenseFile, "license-file", "l", "", "kubemq license file")
	cmd.PersistentFlags().IntVarP(&o.replicas, "replicas", "r", 3, "kubemq cluster replicas")

	return cmd
}

func (o *InstallOptions) Complete(args []string) error {
	if o.licenseFile != "" {
		data, err := os.ReadFile(o.licenseFile)
		if err != nil {
			return fmt.Errorf("error reading license file, %w", err)
		}
		o.licenseData = string(data)
	}
	var err error
	o.kubeCtlPath, err = exec.LookPath("kubectl")
	if err != nil {
		return fmt.Errorf("kubectl not found, %w", err)
	}
	return nil
}

func (o *InstallOptions) Validate() error {
	if o.key == "" && o.licenseData == "" {
		return fmt.Errorf("license key or license file must be provided")
	}

	return nil
}

func (o *InstallOptions) Run(ctx context.Context) error {
	utils.Printlnf("installing KubeMQ cluster %s in namespace %s with %d replicas ...", o.name, o.namespace, o.replicas)
	utils.Printlnf("install KubeMQ CRDs...")
	if err := o.saveToFile("kumemq-init.yaml", []byte(InitTemplate)); err != nil {
		return err
	}
	if err := o.apply("kumemq-init.yaml"); err != nil {
		return err
	}
	utils.Printlnf("install KubeMQ CRDs...done")
	time.Sleep(2 * time.Second)
	utils.Printlnf("install KubeMQ Operator...")
	opTmpl, err := NewOperatorManifest().SetNamespace(o.namespace).Manifest()
	if err != nil {
		return err
	}
	if err := o.saveToFile("kumemq-operator.yaml", []byte(opTmpl)); err != nil {
		return err
	}
	if err := o.apply("kumemq-operator.yaml"); err != nil {
		return err
	}
	utils.Printlnf("install KubeMQ Operator...done")
	time.Sleep(1 * time.Second)
	utils.Printlnf("install KubeMQ Cluster...")
	clusterTmpl, err := NewClusterManifest().SetNamespace(o.namespace).SetName(o.name).SetReplicas(o.replicas).SetKey(o.key).SetLicense(o.licenseData).Manifest()
	if err != nil {
		return err
	}
	if err := o.saveToFile("kumemq-cluster.yaml", []byte(clusterTmpl)); err != nil {
		return err
	}
	if err := o.apply("kumemq-cluster.yaml"); err != nil {
		return err
	}
	utils.Printlnf("install KubeMQ Cluster...done")
	utils.Printlnf("installing KubeMQ cluster %s in namespace %s with %d replicas ,done", o.name, o.namespace, o.replicas)
	utils.Printlnf("Run `kubemqctl get dashboard` to load KubeMQ dashboard in browser")

	return nil
}

func (o *InstallOptions) saveToFile(fileName string, data []byte) error {
	err := os.WriteFile(fileName, data, 0644)
	if err != nil {
		return fmt.Errorf("error saving file %s, %w", fileName, err)
	}
	return nil
}

func (o *InstallOptions) apply(fileName string) error {
	// Set the arguments for the command
	args := []string{"apply", "-f", fileName}

	// Execute the command
	cmd := exec.Command(o.kubeCtlPath, args...)
	// Get the stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error getting stdout pipe, %w", err)
	}
	// Get the stderr pipe
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error getting stderr pipe, %w", err)
	}
	// combine stdout and stderr
	outputReader := io.MultiReader(stdout, stderr)

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error starting command, %w", err)
	}
	scanner := bufio.NewScanner(outputReader)
	for scanner.Scan() {
		utils.Printlnf(scanner.Text())
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
