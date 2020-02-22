## kubemqctl queues list

Get a list of 'queues' channels / clients command

### Synopsis

List command allows to get a list of 'queues' channels / clients with details

```
kubemqctl queues list [flags]
```

### Examples

```

	# Get a list of queues / clients
	kubemqctl queue list
	
	# Get a list of queues / clients filtered by 'some-queue' channel only
	kubemqctl queue list -f some-queue

```

### Options

```
  -f, --filter string   set filter for channel / client name
  -h, --help            help for list
```

### SEE ALSO

* [kubemqctl queues](kubemqctl_queues.md)	 - Execute Kubemq 'queues' commands


