module github.com/kubemq-io/kubemqctl

require (
	github.com/AlecAivazis/survey/v2 v2.2.7
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/color v1.10.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-resty/resty/v2 v2.4.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/websocket v1.4.2
	github.com/json-iterator/go v1.1.10
	github.com/kubemq-hub/builder v0.6.2
	github.com/kubemq-io/kubemq-go v1.4.7
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.5.1
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.17.4
	k8s.io/apiextensions-apiserver v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4

)

go 1.15

//replace github.com/kubemq-hub/builder => ../../kubemq-hub/builder
