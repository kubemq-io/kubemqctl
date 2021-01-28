module github.com/kubemq-io/kubemqctl

require (
	github.com/AlecAivazis/survey/v2 v2.2.7
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fatih/color v1.7.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-resty/resty v1.8.0
	github.com/go-resty/resty/v2 v2.0.0
	github.com/google/uuid v1.1.2
	github.com/gorilla/websocket v1.4.0
	github.com/json-iterator/go v1.1.8
	github.com/kubemq-hub/builder v0.5.9
	github.com/kubemq-io/kubemq-go v1.4.6
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/browser v0.0.0-20180916011732-0a3d74bf9ce4
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.17.4
	k8s.io/apiextensions-apiserver v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4

)

go 1.15

//replace github.com/kubemq-hub/builder => ../../kubemq-hub/builder
