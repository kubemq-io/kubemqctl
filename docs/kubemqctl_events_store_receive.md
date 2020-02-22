## kubemqctl events_store receive

Receive a messages from an 'events store'

### Synopsis

Receive (Subscribe) command allows to consume messages from an 'events store' with options to set offset parameters

```
kubemqctl events_store receive [flags]
```

### Examples

```

	# Receive messages from an 'events store' channel (blocks until next message)
	kubemqctl events_store receive some-channel

	# Receive messages from an 'events channel' with group(blocks until next message)
	kubemqctl events_store receive some-channel -g G1

```

### Options

```
  -g, --group string   set 'events_store' channel consumer group (load balancing)
  -h, --help           help for receive
```

### SEE ALSO

* [kubemqctl events_store](kubemqctl_events_store.md)	 - Execute Kubemq 'events_store' Pub/Sub commands


