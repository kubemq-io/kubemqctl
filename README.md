# KubeTools
Kubetools is a CLI application for testing KubeMQ server or cluster installation.

## Download / Installation
KubeTools executable can be downloaded for 4 operating systems:
1. Windows 64bit
2. Mac OS 64bit
3. Linux 64bit
4. Linux 32bit

## Usage
Run `KubeTools` will prompt the following:
```
Usage:
  kubetools [command]

Available Commands:
  health      Call kubemq health endpoint
  help        Help about any command
  metrics     Call kubemq metrics endpoint
  test        test your kubemq installation
  version     print kubemq version

Flags:
  -h, --help   help for kubetools

Use "kubetools [command] --help" for more information about a command.

```

### Test
Run `kubeTools test` or `kubetools t`

### Health
Run `kubeTools health` or `kubetools h`

### Metrics
Run `kubeTools metrics` or `kubetools m`

## Configuration
KubeTools require `.config.yaml` file for connections variables. Default configuration:

```
healthAddress: "http://localhost:8080/health" # the address of Health endpoint , you can replace the localhost:8080 with your address
metricsAddress: "http://localhost:8080/metrics" #the address of Health endpoint, you can replace the localhost:8080 with your address
connections:
  - kind: 1 # 1 - grpc 2- rest
    host: "localhost" # host destination
    port: 50000 # port destination
    isSecured: false # set using https
    certFile: "" # set location of cert file
  - kind: 2 # 1 - grpc 2- rest
    host: "localhost" # host destination
    port: 9090  # port destination
    isSecured: false  # set using https
    certFile: "" # set location of cert file - not in use for Rest

```
