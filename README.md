# KubeTools
Kubetools is a CLI application for testing KubeMQ server or cluster installation. In addition, Kubetools provides several utilities such as monitoring channel traffic, send messages, and subscribe to channels.

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
  pubsub      send and receive pub/sub real-time and persistent events
  queue       send and receive queue messages
  rpc         send and receive rpc commands and queries
  test        test your kubemq installation
  version     print kubemq version


Flags:
  -h, --help   help for kubetools

Use "kubetools [command] --help" for more information about a command.

```

### Test
Run `kubetools test` or `kubetools t` for running various tests, checking KubeMQ installation and proper configuration.

### Monitor
Run `kubetools mon` or `kubetools m` to enter into a monitoring channels mode.

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
  queue        monitor queue channels

Flags:
  -h, --help   help for mon

Use "kubetools mon [command] --help" for more information about a command.
```

#### Monitor Events Channel
Run `kubetools mon events <ChannelName>` or `kubetools m e <ChannelName>` will monitor and show all traffic in `<ChannelName>` Events channel.

#### Monitor Events Store Channel
Run `kubetools mon events_store <ChannelName>` or `kubetools m es <ChannelName>` will monitor and show all traffic in `<ChannelName>` Events Store channel.

#### Monitor Commands Channel
Run `kubetools mon commands <ChannelName>` or `kubetools m c <ChannelName>` will monitor and show all traffic in `<ChannelName>` Commands channel.

#### Monitor Query Channel
Run `kubetools mon queries <ChannelName>` or `kubetools m q <ChannelName>` will monitor and show all traffic in `<ChannelName>` Queries channel.

#### Monitor Queue Channel
Run `kubetools mon queue <ChannelName>` or `kubetools m qu <ChannelName>` will monitor and show all traffic in `<ChannelName>` Queue channel.


### PubSub
Run `kubetools pubsub` or `kubetools p` for publish and subscribe real-time and persistent events.

Available PubSub commands:

```
Usage:
  kubetools pubsub [command]

Aliases:
  pubsub, p, ps

Available Commands:
  receive     receive pub/sub real-time and persistent events
  send        send pub/sub real-time and persistent events

Flags:
  -h, --help                     help for pubsub
  -t, --pubsubTransport string   set transport type, grpc or rest (default "grpc")

Use "kubetools pubsub [command] --help" for more information about a command.
```


### PubSub Send
Run `kubetools pubsub send` or `kubetools p s` to publish real-time and persistent events.

```
Usage:
  kubetools pubsub send [command]

Aliases:
  send, s

Available Commands:
  events       send real-time event to a channel
  events_store send persistent event to a channel

Flags:
  -h, --help   help for send

Global Flags:
  -t, --pubsubTransport string   set transport type, grpc or rest (default "grpc")

Use "kubetools pubsub send [command] --help" for more information about a command.
```

##### Send Message to an Events Channel
Run `kubetools pubsub send event <ChannelName> <YourMessage>` or `kubetools p s e <ChannelName> <YourMessage>` for sending `<YourMessage>` to Events channel `<ChannelName>`

##### Send Message to an Events Store Channel
Run `kubetools pubsub send events_store <ChannelName> <YourMessage>` or `kubetools p s es <ChannelName> <YourMessage>` for sending `<YourMessage>` to Events Store channel `<ChannelName>`


### PubSub Receive
Run `kubetools pubsub receive` or `kubetools p r` for receiving real-time and persistent events.

```
Usage:
  kubetools pubsub receive [command]

Aliases:
  receive, rec, r

Available Commands:
  events       subscribe to receive real-time events from a channel
  events_store subscribe to receive persistent events from a channel

Flags:
  -h, --help   help for receive

Global Flags:
  -t, --pubsubTransport string   set transport type, grpc or rest (default "grpc")

Use "kubetools pubsub receive [command] --help" for more information about a command.
```

##### Receive PubSub Events From Channel
Run `Kubetools pubsub receive event <ChannelName>` or `Kubetools p r e <ChannelName>` for receiving messages in Events channel `<ChannelName>`.

##### Receive PubSub Events Store From Channel
Run `Kubetools pubsub receive events_store <ChannelName>` or `Kubetools p r es <ChannelName>` for receiving messages in Events Store channel `<ChannelName>`.


### RPC
Run `kubetools rpc` or `kubetools r` for sending and receiving RPC calls of commands and queries.

Available RPC commands:

```
Usage:
  kubetools rpc [command]

Aliases:
  rpc, r

Available Commands:
  receive     receive commands or queries
  send        send commands and queries

Flags:
  -h, --help                  help for rpc
  -t, --rpcTransport string   set transport type, grpc or rest (default "grpc")

Use "kubetools rpc [command] --help" for more information about a command.
```


### RPC Send
Run `kubetools rpc send` or `kubetools r s` for sending command and query rpc calls.

```
Usage:
  kubetools rpc send [command]

Aliases:
  send, s

Available Commands:
  command     send command to a channel
  query       send query to a channel

Flags:
  -h, --help   help for send

Global Flags:
  -t, --rpcTransport string   set transport type, grpc or rest (default "grpc")

Use "kubetools rpc send [command] --help" for more information about a command.
```

##### Send RPC Commands To Channel
Run `kubetools rpc send command <ChannelName> <YourMessage>` or `kubetools r s c <ChannelName> <YourMessage>` for sending `<YourMessage>` to Commands channel `<ChannelName>`

Additional configuration available:

```
Usage:
  kubetools rpc send command [flags]

Aliases:
  command, c

Flags:
  -h, --help                help for command
  -o, --rpcSendTimout int   set command timeout in msec (default 10000)

Global Flags:
  -t, --rpcTransport string   set transport type, grpc or rest (default "grpc")

```

##### Send RPC Query To Channel
Run `kubetools rpc send query <ChannelName> <YourMessage>` or `kubetools r s q <ChannelName> <YourMessage>` for sending `<YourMessage>` to Queries channel `<ChannelName>`

Additional configuration available:

```
Usage:
  kubetools rpc send query [flags]

Aliases:
  query, q

Flags:
  -h, --help                help for query
  -o, --rpcSendTimout int   set query timeout in msec (default 10000)

Global Flags:
  -t, --rpcTransport string   set transport type, grpc or rest (default "grpc")

```


#### RPC RECEIVE
Run `kubetools rpc receive` or `kubetools r r` for receiving commands and queries calls.

```
Usage:
  kubetools rpc receive [command]

Aliases:
  receive, rec, r

Available Commands:
  command     subscribe to receive commands from a channel
  query       subscribe to receive queries from a channel

Flags:
  -h, --help   help for receive

Global Flags:
  -t, --rpcTransport string   set transport type, grpc or rest (default "grpc")

Use "kubetools rpc receive [command] --help" for more information about a command.

```

##### Receive RPC Commands From Channel
Run `Kubetools rpc receive command <ChannelName>` or `Kubetools r r c <ChannelName>` for receiving messages in Commands channel `<ChannelName>` and send an acknowledge back to the sender.

##### Receive RPC Queries From Channel
Run `Kubetools rpc receive query <ChannelName>` or `Kubetools r r q <ChannelName>` for receiving messages in Queries channel `<ChannelName>` and echo back the same message to the sender.


### Queue
Run `kubetools queue` or `kubetools q` for sending and receiving queue messages.

Available Queue commands:

```
Usage:
  kubetools queue [command]

Aliases:
  queue, q

Available Commands:
  ack         acl all messages in a queue
  peak        peak messages from a queue
  receive     receive messages from a queue
  send        send message to a queue

Flags:
  -h, --help                    help for queue
  -t, --queueTransport string   set transport type, grpc or rest (default "grpc")

Use "kubetools queue [command] --help" for more information about a command.
```


#### Send To a Queue
Run `kubetools queue send` or `kubetools q s` for sending queue messages.

```
Usage:
  kubetools queue send [flags]

Aliases:
  send, s

Flags:
  -h, --help                 help for send
  -d, --sendDelay int        set queue message send delay seconds
  -e, --sendExpiration int   set queue message expiration seconds

Global Flags:
  -t, --queueTransport string   set transport type, grpc or rest (default "grpc")
```

#### Receive From a Queue
Run `kubetools queue receive` or `kubetools q r` for receiving queue messages.

```
Usage:
  kubetools queue receive [flags]

Aliases:
  receive, r

Flags:
  -h, --help                  help for receive
  -i, --receiveMessages int   set how many messages we want to get from queue (default 1)
  -w, --receiveWait int       set how many seconds to wait for queue messages (default 2)

Global Flags:
  -t, --queueTransport string   set transport type, grpc or rest (default "grpc")

```

#### Peak into a Queue
Run `kubetools queue peak` or `kubetools q p` for peaking queue messages.

```
Usage:
  kubetools queue peak [flags]

Aliases:
  peak, p

Flags:
  -h, --help                  help for peak
  -i, --receiveMessages int   set how many messages we peak to get from queue (default 1)
  -w, --receiveWait int       set how many seconds to peak for queue messages (default 2)

Global Flags:
  -t, --queueTransport string   set transport type, grpc or rest (default "grpc")
```

#### Ack All Messages in a Queue
Run `kubetools queue ack` or `kubetools q a` for ack all queue messages.

```
Usage:
  kubetools queue ack [flags]

Aliases:
  ack, a

Flags:
  -h, --help              help for ack
  -w, --receiveWait int   set how many seconds wait to ack all messages in queue (default 2)

Global Flags:
  -t, --queueTransport string   set transport type, grpc or rest (default "grpc")
```

### Health
Run `kubetools health` or `kubetools h`

### Metrics
Run `kubetools metrics` or `kubetools m`



## Configuration
KubeTools require `.config.yaml` File for connections variables. Default configuration:

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
