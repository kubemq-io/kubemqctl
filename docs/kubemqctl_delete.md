## kubemqctl delete

Executes delete commands

### Synopsis

Executes delete commands

```
kubemqctl delete [flags]
```

### Examples

```

	# Execute delete Kubemq cluster
	kubemqctl delete cluster
	
	# Execute delete Kubemq Operator
	kubemqctl delete operator	

   # Execute delete Kubemq Connector
	kubemqctl delete connector

   # Execute delete Kubemq Components
	kubemqctl delete connector

```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
      --config string   set kubemqctl configuration file (default "./.kubemqctl.yaml")
```

### SEE ALSO

* [kubemqctl](kubemqctl.md)	 - 
* [kubemqctl delete cluster](kubemqctl_delete_cluster.md)	 - Delete Kubemq cluster
* [kubemqctl delete components](kubemqctl_delete_components.md)	 - Delete Kubemq components
* [kubemqctl delete connector](kubemqctl_delete_connector.md)	 - Delete Kubemq connector
* [kubemqctl delete operator](kubemqctl_delete_operator.md)	 - Delete Kubemq operator

###### Auto generated by spf13/cobra on 10-Jul-2023
