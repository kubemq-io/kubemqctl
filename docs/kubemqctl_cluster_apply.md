## kubemqctl cluster apply

Apply a Kubemq cluster command

### Synopsis

Apply command allows an update to a Kubemq StatefulSet configuration with a yaml file

```
kubemqctl cluster apply [flags]
```

### Examples

```

	# Apply Kubemq cluster deployment
	# kubemqctl cluster apply kubemq-cluster.yaml 

	# Apply Kubemq cluster deployment with watching status and events
	# kubemqctl cluster apply kubemq-cluster.yaml -w -s


```

### Options

```
  -f, --file string   set yaml configuration file
  -h, --help          help for apply
  -s, --status        watch and print apply StatefulSet status
  -w, --watch         watch and print apply StatefulSet events
```

### SEE ALSO

* [kubemqctl cluster](kubemqctl_cluster.md)	 - Executes Kubemq cluster management commands


