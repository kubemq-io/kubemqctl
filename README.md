# KubeTools
Kubetools is a CLI application for testing KubeMQ server or cluster installation. In addition Kubetools provides several utilities such monitoring channel traffic, send messages and subscribe to channels.

## Download / Installation
KubeTools executable can be downloaded from bin folder for 4 operating systems:
1. Windows 64bit
2. Mac OS 64bit
3. Linux 64bit
4. Linux 32bit

## Usage
Run `Kubetools` will prompt the following:
```
Usage:
  kubetools [command]

Available Commands:
  health      Call kubemq health endpoint
  help        Help about any command
  metrics     Call kubemq metrics endpoint
  mon         monitor messages/requests channels
  send        send event / event_store / command / query
  subscribe   subscribe to events / events_store / commands / queries
  test        test your kubemq installation
  version     print kubemq version

Flags:
  -h, --help   help for kubetools

Use "kubetools [command] --help" for more information about a command.

```

### Test
Run `kubetools test` or `kubetools t` for running various tests ,checking KubeMQ installation and proper configuration.

### Monitor
Run `kubetools mon` or `kubetools m` to enter monitoring channels mode.

Available sub commands:

```
Usage:
  kubetools mon [command]

Aliases:
  mon, m

Available Commands:
  commands     monitor commands channels
  events       monitor events channels
  events_store monitor events store channels
  queries      monitor query channels

Flags:
  -h, --help   help for mon

Use "kubetools mon [command] --help" for more information about a command.
```

#### Monitor Events Channel
Run `kubetools mon events <ChannelName>` or `kubetools m e <ChannelName>` will monitor and show all traffic in <ChannelName> Events channel.

#### Monitor Events Store Channel
Run `kubetools mon events_store <ChannelName>` or `kubetools m es <ChannelName>` will monitor and show all traffic in <ChannelName> Events Store channel.

#### Monitor Commands Channel
Run `kubetools mon commands <ChannelName>` or `kubetools m c <ChannelName>` will monitor and show all traffic in <ChannelName> Commands channel.

#### Monitor Query Channel
Run `kubetools mon queries <ChannelName>` or `kubetools m q <ChannelName>` will monitor and show all traffic in <ChannelName> Queries channel.

### Send
Run `kubetools send` or `kubetools s` for sending messages to any channel.

Available sub commands:

```
Usage:
  kubetools send [command]

Aliases:
  send, s

Available Commands:
  command     send command to a channel
  event       send event to a channel
  event_store send event_store to a channel
  query       send query to a channel

Flags:
  -h, --help                   help for send
  -t, --sendTransport string   set transport type, grpc or rest (default "grpc")

Use "kubetools send [command] --help" for more information about a command.

```

#### Send Message to an Events Channel
Run `kubetools send event <ChannelName> <YourMessage>` or `kubetools s e <ChannelName> <YourMessage>` for sending <YourMessage> to Events channel <ChannelName>

#### Send Message to an Events Store Channel
Run `kubetools send events_store <ChannelName> <YourMessage>` or `kubetools s es <ChannelName> <YourMessage>` for sending <YourMessage> to Events Store channel <ChannelName>

#### Send Message to a Commands Channel
Run `kubetools send command <ChannelName> <YourMessage>` or `kubetools s c <ChannelName> <YourMessage>` for sending <YourMessage> to Commands channel <ChannelName>

Additional configuration available:

```
Usage:
  kubetools send command [flags]

Aliases:
  command, c

Flags:
  -h, --help             help for command
  -o, --sendTimout int   set command timeout in MSec (default 10000)

Global Flags:
  -t, --sendTransport string   set transport type, grpc or rest (default "grpc")

```

#### Send Message to a Query Channel
Run `kubetools send query <ChannelName> <YourMessage>` or `kubetools s q <ChannelName> <YourMessage>` for sending <YourMessage> to Queries channel <ChannelName>

Additional configuration available:

```
Usage:
  kubetools send query [flags]

Aliases:
  query, q

Flags:
  -h, --help             help for query
  -o, --sendTimout int   set query timeout in MSec (default 10000)

Global Flags:
  -t, --sendTransport string   set transport type, grpc or rest (default "grpc")

```

### Subscribe
Run `kubetools send` or `kubetools s` for subscribing to any channel ,show messages received and echo messages back in Commands and Queries channels.

```
Usage:
  kubetools subscribe [command]

Aliases:
  subscribe, sub

Available Commands:
  command     subscribe to a command to a channel
  event       subscribe to an events channel
  event_store subscribe to an event_store channel
  query       subscribe to a query channel

Flags:
  -h, --help                        help for subscribe
  -g, --subscribeGroup string       set optional group for a channel
  -t, --subscribeTransport string   set transport type, grpc or rest (default "grpc")


```

#### Subscribe to Events Channel
Run `Kubetools subscribe event <ChannelName>` or `Kubetools sub e <ChannelName>` for receiving messages in Events channel <ChannelName>.

#### Subscribe to Events Store Channel
Run `Kubetools subscribe event_store <ChannelName>` or `Kubetools sub es <ChannelName>` for receiving messages in Events Store channel <ChannelName>.

#### Subscribe to Commands Channel
Run `Kubetools subscribe command <ChannelName>` or `Kubetools sub c <ChannelName>` for receiving messages in Commands channel <ChannelName> and send acknowledge back to the sender.

#### Subscribe to Queries Channel
Run `Kubetools subscribe query <ChannelName>` or `Kubetools sub q <ChannelName>` for receiving messages in Queries channel <ChannelName> and echo back the same message to the sender.

### Health
Run `kubetools health` or `kubetools h`

### Metrics
Run `kubetools metrics` or `kubetools m`



## Configuration
KubeTools require `.config.yaml` file for connections variables. Default configuration:

```
healthAddress: "http://localhost:8080/health" # the address of Health endpoint , you can replace the localhost:8080 with your address
metricsAddress: "http://localhost:8080/metrics" #the address of Health endpoint, you can replace the localhost:8080 with your address
monitorAddress: "ws://localhost:8080/v1/stats" #the address of Monitor endpoint, you can replace the localhost:8080 with your address
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
