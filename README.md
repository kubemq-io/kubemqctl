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

```bash
curl -L https://github.com/kubemq-io/kubemqctl/releases/download/latest/kubemqctl_darwin_amd64 -o /usr/local/bin/kubemqctl 
chmod +x /usr/local/bin/kubemqctl
```

### Linux 64 bits:

```bash
curl -L https://github.com/kubemq-io/kubemqctl/releases/download/latest/kubemqctl_linux_amd64 -o /usr/local/bin/kubemqctl
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


## Documantation
Please visit our [docs](https://docs.kubemq.io/kubemqctl/kubemqctl.html) for detailed Kubemqctl documentation.


## Support
if you encounter any issues, please open an issue here,
In addition, you can reach us for support by:
- [**Email**](mailto://support@kubemq.io)
- [**Slack**](https://kubmq.slack.com)
