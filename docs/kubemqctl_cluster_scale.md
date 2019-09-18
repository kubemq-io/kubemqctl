## kubemqctl cluster scale

Scale KubeMQ cluster replicas command

### Synopsis

Scale command allows ro scale KubeMQ cluster replicas

```
kubemqctl cluster scale [flags]
```

### Examples

```

	# Scale KubeMQ cluster StatefulSet 
	kubemqctl cluster cluster scale 5

	# Scale KubeMQ cluster StatefulSet with streaming real-time events and status
	kubemqctl cluster scale -w -s 

	# Scale KubeMQ cluster StatefulSet to 0
	kubemqctl cluster scale 0

```

### Options

```
  -h, --help     help for scale
  -s, --status   watch and print Scale StatefulSet status
  -w, --watch    watch and print Scale StatefulSet events
```

### SEE ALSO

* [kubemqctl cluster](kubemqctl_cluster.md)	 - Executes KubeMQ cluster management commands


