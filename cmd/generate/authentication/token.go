package authentication

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type TokenOptions struct {
	cfg    *config.Config
	verify bool
}

var tokenExamples = `
	# Execute generate authentication JWT token
 	kubemqctl generate auth token

	# Execute JWT token verification
 	kubemqctl generate auth token -v
`
var tokenLong = `Generate and validate JWT tokens`
var tokenShort = `Generate and validate JWT tokens`

func NewCmdToken(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &TokenOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "token",
		Aliases: []string{"t"},
		Short:   tokenShort,
		Long:    tokenLong,
		Example: tokenExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().BoolVarP(&o.verify, "verify", "v", false, "set to verify a token")
	return cmd
}

func (o *TokenOptions) Complete(args []string, transport string) error {
	return nil
}

func (o *TokenOptions) Validate() error {
	return nil
}
func (o *TokenOptions) Run(ctx context.Context) error {
	if o.verify {
		return o.verifyToken()
	}
	return o.generateToken()
}
func (o *TokenOptions) verifyToken() error {
	publicKey, signatureType, err := getPublicKey()
	if err != nil {
		return err
	}

	verifySigner, err := CreateVerifySignature(publicKey, signatureType)
	if err != nil {
		return err
	}
	token, err := o.getToken()
	if err != nil {
		return err
	}

	type jwtResponse struct {
		jwt.StandardClaims
	}
	response := &jwtResponse{}
	err = verifySigner.Verify(token, response)
	if err != nil {
		return err
	}
	utils.Println("Token Verified:")
	utils.Printlnf("Id: %s", response.Id)
	utils.Printlnf("Issued to: %s", response.Audience)
	utils.Printlnf("Issued at: %s", time.Unix(response.IssuedAt, 0).Format("2006-01-02 15:04:05"))
	utils.Printlnf("Expired at: %s", time.Unix(response.ExpiresAt, 0).Format("2006-01-02 15:04:05"))

	return nil
}

func (o *TokenOptions) generateToken() error {

	privateKey, signatureType, err := getPrivateKey()
	if err != nil {
		return err
	}
	signer, err := CreateSignSignature(privateKey, signatureType)
	if err != nil {
		return err
	}
	claims, err := getClaims()
	if err != nil {
		return err
	}

	for i, claim := range claims {
		token, err := signer.Sign(claim)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s.key", claim.Id), []byte(token), 0600)
		if err != nil {
			return err
		} else {
			utils.Printlnf("JWT Token #%d exported to file %s.key.", i+1, claim.Id)
		}

	}
	return nil
}
func (o *TokenOptions) getToken() (string, error) {
	token := ""
	promptGetToken := &survey.Editor{
		Message: "Copy & Paste Authentication Token data",
		Default: "",
		Help:    "Set Authentication Token data",
	}
	err := survey.AskOne(promptGetToken, &token, survey.WithValidator(survey.MinLength(1)))
	if err != nil {
		return "", err
	}

	cleaned := strings.Replace(token, " ", "", -1)
	cleaned = strings.Replace(cleaned, "\t", "", -1)
	cleaned = strings.Replace(cleaned, "\r", "", -1)
	cleaned = strings.Replace(cleaned, "\n", "", -1)

	return cleaned, nil
}
func getClaims() ([]jwt.StandardClaims, error) {
	issuedTo := ""
	promptIssuedTo := &survey.Input{
		Message: "Issue JWT Token to ?",
		Default: "",
		Help:    "Set the name of the token owner",
	}
	err := survey.AskOne(promptIssuedTo, &issuedTo)
	if err != nil {
		return nil, err
	}

	expireAt := &TimeAnswer{}
	promptExpireAt := &survey.Input{
		Message: "Set JWT Token Expiration time, i.e 1h or 2022-01-02:",
		Default: "24h",
		Help:    "Set JWT token expiration time in duration or time formats",
	}
	err = survey.AskOne(promptExpireAt, expireAt, survey.WithValidator(expireAt.Validate))
	if err != nil {
		return nil, err
	}
	numOfTokens := "1"
	numOfTokensPrompt := &survey.Input{
		Message: "Set number of JWT Tokens to generate",
		Default: "1",
		Help:    "Set number of JWT Tokens to generate",
		Suggest: nil,
	}
	err = survey.AskOne(numOfTokensPrompt, &numOfTokens, survey.WithValidator(
		func(ans interface{}) error {
			val := ans.(string)
			_, err := strconv.Atoi(val)
			if err != nil {
				return err
			}
			return nil
		}))
	if err != nil {
		return nil, err
	}
	n, err := strconv.Atoi(numOfTokens)
	if err != nil {
		return nil, err
	}
	if n < 1 {
		return nil, fmt.Errorf("number of jwt token must have at least 1")
	}
	var list []jwt.StandardClaims
	for i := 0; i < n; i++ {
		list = append(list, jwt.StandardClaims{
			Audience:  issuedTo,
			ExpiresAt: expireAt.value.Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "",
			NotBefore: 0,
			Subject:   "JWT Token",
		})
	}
	return list, err
}

func getPrivateKey() ([]byte, string, error) {
	privateKey := ""
	promptPrivateKey := &survey.Editor{
		Message:  "Copy & Paste Private Key",
		FileName: "*.pem",
	}

	err := survey.AskOne(promptPrivateKey, &privateKey, survey.WithValidator(survey.MinLength(1)))
	if err != nil {
		return nil, "", err
	}

	signatureType := ""
	signatureTypePrompt := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  "Select Private Key signature type:",
		Options:  []string{"HS256", "HS384", "HS512", "RS256", "RS384", "RS512", "ES256", "ES384", "ES512"},
		Default:  "RS512",
		Help:     "Select Private Key signature type",
	}
	err = survey.AskOne(signatureTypePrompt, &signatureType)
	if err != nil {
		return nil, "", err
	}
	return []byte(privateKey), signatureType, nil
}
func getPublicKey() ([]byte, string, error) {
	publicKey := ""
	promptPublicKey := &survey.Editor{
		Message:  "Copy & Paste Public Key",
		FileName: "*.pem",
	}

	err := survey.AskOne(promptPublicKey, &publicKey, survey.WithValidator(survey.MinLength(1)))
	if err != nil {
		return nil, "", err
	}
	signatureType := ""
	signatureTypePrompt := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  "Select Public Key signature type:",
		Options:  []string{"HS256", "HS384", "HS512", "RS256", "RS384", "RS512", "ES256", "ES384", "ES512"},
		Default:  "RS512",
		Help:     "Select Public Key signature type",
	}
	err = survey.AskOne(signatureTypePrompt, &signatureType)
	if err != nil {
		return nil, "", err
	}
	return []byte(publicKey), signatureType, nil
}
