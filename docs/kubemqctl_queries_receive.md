## kubemqctl queries receive

Receive a message from a 'queries' channel

### Synopsis

Receive (Subscribe) command allows to receive a message from a 'queries' channel and response with appropriate reply

```
kubemqctl queries receive [flags]
```

### Examples

```

	# Receive 'queries'  from a 'queries' channel (blocks until next message)
	kubemqctl queries receive some-channel

	# Receive 'queries' from a 'queries' channel with group(blocks until next message)
	kubemqctl queries receive some-channel -g G1

```

### Options

```
  -a, --auto-response   set auto response executed query
  -g, --group string    set 'queries' channel consumer group (load balancing)
  -h, --help            help for receive
```

### SEE ALSO

* [kubemqctl queries](kubemqctl_queries.md)	 - Execute Kubemq 'queries' RPC based commands


