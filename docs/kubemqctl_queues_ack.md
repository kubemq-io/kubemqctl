## kubemqctl queues ack

Ack all messages in a 'queues' channel

### Synopsis

Ack command allows to clear / remove / ack all messages in a 'queues' channel

```
kubemqctl queues ack [flags]
```

### Examples

```

	# Ack all messages in a 'queues' channel 'some-channel' with 2 seconds of wait to complete operation
	kubemqctl queue ack some-channel
	
	# Ack all messages in a 'queues' channel 'some-long-queue' with 30 seconds of wait to complete operation
	kubemqctl queue ack some-long-queue -w 30

```

### Options

```
  -h, --help       help for ack
  -w, --wait int   set how many seconds to wait for ack all messages (default 2)
```

### SEE ALSO

* [kubemqctl queues](kubemqctl_queues.md)	 - Execute KubeMQ 'queues' commands


