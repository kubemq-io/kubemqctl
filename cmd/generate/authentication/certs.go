package authentication

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

type CertsOptions struct {
	cfg    *config.Config
	output string
}

var certsExamples = `
	# Execute generate authentication rsa certificates
 	kubemqctl generate auth certs
`
var certsLong = `Generate JWT certificates`
var certsShort = `Generate JWT certificates`

func NewCmdCerts(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &CertsOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "certs",
		Aliases: []string{"cert"},
		Short:   certsShort,
		Long:    certsLong,
		Example: certsExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().StringVarP(&o.output, "output", "o", "", "set output filename")
	return cmd
}

func (o *CertsOptions) Complete(args []string, transport string) error {
	return nil
}

func (o *CertsOptions) Validate() error {
	return nil
}

func (o *CertsOptions) Run(ctx context.Context) error {
	return generateCerts(o.output)
}

func generateCerts(filePrefix string) error {
	var isExport bool
	var privateKeyFileName string
	var publicKeyFileName string
	if filePrefix != "" {
		isExport = true
		privateKeyFileName = fmt.Sprintf("%s-private.pem", filePrefix)
		publicKeyFileName = fmt.Sprintf("%s-public.pem", filePrefix)
	}
	// generate key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	publicKey := &privateKey.PublicKey
	// dump private key to file
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	if isExport {

		privateWriter, err := os.Create(privateKeyFileName)
		if err != nil {
			return err
		}
		err = pem.Encode(privateWriter, privateKeyBlock)
		if err != nil {
			return err
		}
		utils.Printlnf("RSA Private key exported to %s file", privateKeyFileName)
	} else {
		utils.Println("RSA Private key:")
		err = pem.Encode(os.Stdout, privateKeyBlock)
		if err != nil {
			return err
		}
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	if isExport {

		publicWriter, err := os.Create(publicKeyFileName)
		if err != nil {
			return err
		}
		err = pem.Encode(publicWriter, publicKeyBlock)
		if err != nil {
			return err
		}
		utils.Printlnf("RSA Public key exported to %s file", publicKeyFileName)
	} else {
		utils.Println("RSA Public key:")
		err = pem.Encode(os.Stdout, publicKeyBlock)
		if err != nil {
			return err
		}
	}

	return nil
}
