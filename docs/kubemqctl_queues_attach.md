## kubemqctl queues attach

Attach to 'queues' channels command

### Synopsis

Attach command allows to display 'queues' channel content for debugging proposes

```
kubemqctl queues attach [flags]
```

### Examples

```

	# Attach to all active 'queues' channels and output running messages
	kubemqctl queue attach all
	
	# Attach to some-queue queue channel and output running messages
	kubemqctl queue attach some-queue

	# Attach to some-queue1 and some-queue2 queue channels and output running messages
	kubemqctl queue attach some-queue1 some-queue2 

	# Attach to some-queue queue channel and output running messages filter by include regex (some*)
	kubemqctl queue attach some-queue -i some*

	# Attach to some-queue queue channel and output running messages filter by exclude regex (not-some*)
	kubemqctl queue attach some-queue -e not-some*

```

### Options

```
  -e, --exclude stringArray   set (regex) strings to exclude
  -h, --help                  help for attach
  -i, --include stringArray   aet (regex) strings to include
```

### SEE ALSO

* [kubemqctl queues](kubemqctl_queues.md)	 - Execute Kubemq 'queues' commands


