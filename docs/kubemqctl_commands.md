## kubemqctl commands

Execute KubeMQ 'commands' RPC commands

### Synopsis

Execute KubeMQ 'commands' RPC commands

```
kubemqctl commands [flags]
```

### Examples

```

	# Execute send commands 
	kubemqctl commands send

	# Execute receive commands
	kubemqctl commands receive

	# Execute attach to 'commands' channel
	kubemqctl commands attach

```

### Options

```
  -h, --help   help for commands
```

### SEE ALSO

* [kubemqctl](kubemqctl.md)	 - 
* [kubemqctl commands attach](kubemqctl_commands_attach.md)	 - Attach to 'commands' channels command
* [kubemqctl commands receive](kubemqctl_commands_receive.md)	 - Receive a message from 'commands' channel command
* [kubemqctl commands send](kubemqctl_commands_send.md)	 - Send messages to 'commands' channel command


