## kubemqctl events_store

Execute Kubemq 'events_store' Pub/Sub commands

### Synopsis

Execute Kubemq 'events_store' Pub/Sub commands

```
kubemqctl events_store [flags]
```

### Examples

```

	# Execute send 'events store' command 
	kubemqctl events_store send

	# Execute receive 'events store' command
	kubemqctl events_store receive

	# Execute attach to 'events store' command
	 kubemqctl events_store attach

	# Execute list of 'events store' channels command
 	kubemqctl events_store list

```

### Options

```
  -h, --help   help for events_store
```

### SEE ALSO

* [kubemqctl](kubemqctl.md)	 - 
* [kubemqctl events_store attach](kubemqctl_events_store_attach.md)	 - Attach to events store channels command
* [kubemqctl events_store list](kubemqctl_events_store_list.md)	 - Get a list of 'events store' channels / clients command
* [kubemqctl events_store receive](kubemqctl_events_store_receive.md)	 - Receive a messages from an 'events store'
* [kubemqctl events_store send](kubemqctl_events_store_send.md)	 - Send messages to an 'events store' channel command


