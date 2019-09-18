## kubemqctl cluster create

Create a KubeMQ cluster command

### Synopsis

Create command allows to deploy a KubeMQ cluster with configuration options

```
kubemqctl cluster create [flags]
```

### Examples

```

	# Create default KubeMQ cluster
	# kubemqctl cluster create -t b33600cc-93ef-4395-bba3-13131eb27d5e 

	# Create default KubeMQ cluster and watch events and status
	# kubemqctl cluster create -t b3d3600cc-93ef-4395-bba3-13131eb27d5e -w -s

	# Import KubeMQ cluster yaml file  
	# kubemqctl cluster create -f kubemq-cluster.yaml

	# Create KubeMQ cluster with options
	# kubemqctl cluster create -t b33d30scc-93ef-43565-bba3-13131sb2785e -o

	# Export KubeMQ cluster yaml file    
	# kubemqctl cluster create -t b3d330scc-93qf-4395-bba3-13131sb2785e -e 

```

### Options

```
  -e, --export         generate yaml configuration file output (exporting)
  -f, --file string    import configuration yaml file
  -h, --help           help for create
  -o, --options        create KubeMQ cluster with options
  -s, --status         stream real-time status events during KubeMQ cluster Create command
  -t, --token string   set KubeMQ Token
  -w, --watch          stream real-time events during KubeMQ cluster Create command
```

### SEE ALSO

* [kubemqctl cluster](kubemqctl_cluster.md)	 - Executes KubeMQ cluster management commands


