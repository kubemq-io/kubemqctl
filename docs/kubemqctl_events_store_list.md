## kubemqctl events_store list

Get a list of 'events store' channels / clients command

### Synopsis

List command allows to get a list of 'events store' channels / clients with details

```
kubemqctl events_store list [flags]
```

### Examples

```

	# Get a list of events store channels
	kubemqctl events_store list
	
	# Get a list of events stores channels/ clients filtered by 'some-events-store' channel only
	kubemqctl events_store list -f some-events-store

```

### Options

```
  -f, --filter string   set filter for channel / client name
  -h, --help            help for list
```

### SEE ALSO

* [kubemqctl events_store](kubemqctl_events_store.md)	 - Execute KubeMQ 'events_store' Pub/Sub commands


