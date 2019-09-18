# Kubemqctl

Kubemqctl is a command line interface for [KubeMQ](https://kubemq.io) Kubernetes Message Broker.

```bash
Usage:
  kubemqctl [command]

Available Commands:
  cluster      Executes KubeMQ cluster management commands
  commands     Execute KubeMQ 'commands' RPC commands
  config       Run Kubemqctl configuration wizard command
  events       Execute KubeMQ 'events' Pub/Sub commands
  events_store Execute KubeMQ 'events_store' Pub/Sub commands
  help         Help about any command
  queries      Execute KubeMQ 'queries' RPC based commands
  queues       Execute KubeMQ 'queues' commands

Flags:
  -h, --help      help for kubemqctl
      --version   version for kubemqctl

Use "kubemqctl [command] --help" for more information about a command.

```
## Installation

### Mac OS:

#### Homebrew

Use Homebrew to install 'kubemqctl' on macOS:
``` bash
brew install kubemq-io/homebrew-tap/kubemqctl
```

#### Download Release

Download, extract kubemqctl release, move to local bin and set permissions.

```bash
curl -L https://github.com/kubemq-io/kubemqctl/releases/download/v<version>/kubemqctl_<version>_darwin_amd64.tar.gz | tar -xzv
sudo mv ~/kubemqctl /usr/local/bin
chmod +x /usr/local/bin/kubemqctl
```

### Linux 64 bits:

#### Download Release

Download, extract kubemqctl release, move to local bin and set permissions.

```bash
curl -L https://github.com/kubemq-io/kubemqctl/releases/download/v<version>/kubemqctl_<version>_linux_amd64.tar.gz | tar -xzv
sudo mv ~/kubemqctl /usr/local/bin
chmod +x /usr/local/bin/kubemqctl
```

### Linux 32 bits:

#### Download Release

Download, extract kubemqctl release, move to local bin and set permissions.

```bash
curl -L https://github.com/kubemq-io/kubemqctl/releases/download/v<version>/kubemqctl_<version>_linux_386.tar.gz | tar -xzv
sudo mv ~/kubemqctl /usr/local/bin
chmod +x /usr/local/bin/kubemqctl
```


### Windows:

- [Download the kubemqctl zip file](https://github.com/kubemq-io/kubemqctl/releases/download/v<version>/kubemqctl_<version>_windows_amd64.zip).
- unzip the downloaded file
- Add that directory to your system path to access it from any command prompt
 ([Windows users can follow How to: Add Tool Locations to the PATH Environment Variable in order to add kubemqctl to their PATH](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx)).

## Documantation
Please visit our [docs](https://docs.kubemq.io/tutorials/kubemqctl.html) for detailed Kubemqctl documentation.


## Support
if you encounter any issues, please open an issue here,
In addition, you can reach us for support by:
- [**Email**](mailto://support@kubemq.io)
- [**Slack**](https://kubmq.slack.com)
