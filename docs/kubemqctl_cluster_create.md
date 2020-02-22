## kubemqctl cluster create

Create a Kubemq cluster command

### Synopsis

Create command allows to deploy a Kubemq cluster with configuration options

```
kubemqctl cluster create [flags]
```

### Examples

```

	# Create default Kubemq cluster
	# kubemqctl cluster create -t b33600cc-93ef-4395-bba3-13131eb27d5e 

	# Create default Kubemq cluster and watch events and status
	# kubemqctl cluster create -t b3d3600cc-93ef-4395-bba3-13131eb27d5e -w -s

	# Import Kubemq cluster yaml file  
	# kubemqctl cluster create -f kubemq-cluster.yaml

	# Create Kubemq cluster with options
	# kubemqctl cluster create -t b33d30scc-93ef-43565-bba3-13131sb2785e -o

	# Export Kubemq cluster yaml file    
	# kubemqctl cluster create -t b3d330scc-93qf-4395-bba3-13131sb2785e -e 

```

### Options

```
  -e, --export         generate yaml configuration file output (exporting)
  -f, --file string    import configuration yaml file
  -h, --help           help for create
  -o, --options        create Kubemq cluster with options
  -s, --status         stream real-time status events during Kubemq cluster Create command
  -t, --token string   set Kubemq Token
  -w, --watch          stream real-time events during Kubemq cluster Create command
```

### SEE ALSO

* [kubemqctl cluster](kubemqctl_cluster.md)	 - Executes Kubemq cluster management commands


