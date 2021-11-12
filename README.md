# Kubemqctl

Kubemqctl is a command line interface (CLI) for [Kubemq](https://kubemq.io) [Kubernetes](https://kubernetes.io/) Message Broker.

```bash
Usage:
  kubemqctl [command]

Available Commands:
  commands     Execute Kubemq 'commands' RPC commands
  config       Run Kubemqctl configuration wizard command
  create       Executes Kubemq create commands
  delete       Executes delete commands
  events       Execute Kubemq 'events' Pub/Sub commands
  events_store Execute Kubemq 'events_store' Pub/Sub commands
  generate     Generate various kubemq related artifacts
  get          Executes Kubemq get commands
  help         Help about any command
  manage       Executes Kubemq manage command
  queries      Execute Kubemq 'queries' RPC based commands
  queues       Execute Kubemq 'queues' commands
  scale        Executes Kubemq scale commands
  set          Executes set commands

Flags:
      --config string   set kubemqctl configuration file (default "./.kubemqctl.yaml")
  -h, --help            help for kubemqctl
  -v, --version         version for kubemqctl

Use "kubemqctl [command] --help" for more information about a command.
```
## Installation

### Mac OS:

```bash
curl -L https://github.com/kubemq-io/kubemqctl/releases/download/latest/kubemqctl_darwin_amd64 -o /usr/local/bin/kubemqctl 
chmod +x /usr/local/bin/kubemqctl

```

### Linux 64 bits:

```bash
curl -L https://github.com/kubemq-io/kubemqctl/releases/download/latest/kubemqctl_linux_amd64 -o /usr/local/bin/kubemqctl
chmod +x /usr/local/bin/kubemqctl

```


### Linux 32 bits:

```bash
curl -L https://github.com/kubemq-io/kubemqctl/releases/download/latest/kubemqctl_linux_386 -o /usr/local/bin/kubemqctl
chmod +x /usr/local/bin/kubemqctl

```

### Windows:

Run in PowerShell as administrator:

```powershell
New-Item -ItemType Directory 'C:\Program Files\Kubemqctl'
Invoke-WebRequest https://github.com/kubemq-io/kubemqctl/releases/download/latest/kubemqctl.exe -OutFile 'C:\Program Files\Kubemqctl\kubemqctl.exe'
[Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', [EnvironmentVariableTarget]::Machine) + ';C:\Program Files\Kubemqctl', [EnvironmentVariableTarget]::Machine)
$env:Path += ';C:\Program Files\Kubemqctl'
```

Or manually:

- [Download the latest kubemqctl.exe](https://github.com/kubemq-io/kubemqctl/releases/download/latest/kubemqctl.exe).
- Place the file under e.g. `C:\Program Files\Kubemqctl\kubemqctl.exe`
- Add that directory to your system path to access it from any command prompt

## Kubemq Token

Please visit [Register/Login](https://account.kubemq.io/login/register) to obtain Kubemq token.


## Documantation
Please visit our [docs](https://docs.kubemq.io/kubemqctl/kubemqctl.html) for detailed Kubemqctl documentation.


## Support
if you encounter any issues, please open an issue here,
In addition, you can reach us for support by:
- [**Email**](mailto://support@kubemq.io)
- [**Slack**](https://kubmq.slack.com)
