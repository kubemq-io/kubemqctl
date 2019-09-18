## kubemqctl events receive

Receive a message from 'events' channel command

### Synopsis

Receive (Subscribe) command allows to consume one or many messages from 'events' channel

```
kubemqctl events receive [flags]
```

### Examples

```

	# Receive messages from an 'events' channel (blocks until next message)
	kubemqctl events receive some-channel

	# Receive messages from an 'events' channel with group (blocks until next message)
	kubemqctl events receive some-channel -g G1


```

### Options

```
  -g, --group string   set 'events' channel consumer group (load balancing)
  -h, --help           help for receive
```

### SEE ALSO

* [kubemqctl events](kubemqctl_events.md)	 - Execute KubeMQ 'events' Pub/Sub commands


