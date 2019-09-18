## kubemqctl cluster

Executes KubeMQ cluster management commands

### Synopsis

Executes KubeMQ cluster management commands

```
kubemqctl cluster [flags]
```

### Examples

```

	# Execute create KubeMQ cluster command
	kubemqctl cluster create

	# Execute delete KubeMQ cluster command
	kubemqctl cluster delete

	# Execute describe KubeMQ cluster command
	kubemqctl cluster describe

	# Execute apply KubeMQ cluster command
	kubemqctl cluster apply

	# Execute show KubeMQ cluster logs command
	kubemqctl cluster logs

	# Execute scale KubeMQ cluster command
	kubemqctl cluster scale

	# Execute list of KubeMQ clusters command
	kubemqctl cluster list

	# Execute proxy ports of KubeMQ cluster command
	kubemqctl cluster proxy

	# Execute switch Kubernetes context command
	kubemqctl cluster context

	# Execute cluster web-based dashboard
	kubemqctl cluster dashboard

	# Show cluster events
	kubemqctl cluster events

```

### Options

```
  -h, --help   help for cluster
```

### SEE ALSO

* [kubemqctl](kubemqctl.md)	 - 
* [kubemqctl cluster apply](kubemqctl_cluster_apply.md)	 - Apply a KubeMQ cluster command
* [kubemqctl cluster context](kubemqctl_cluster_context.md)	 - Select kubernetes cluster context command
* [kubemqctl cluster create](kubemqctl_cluster_create.md)	 - Create a KubeMQ cluster command
* [kubemqctl cluster dashboard](kubemqctl_cluster_dashboard.md)	 - Dashboard command allows to start a web view of KubeMQ cluster dashboard
* [kubemqctl cluster delete](kubemqctl_cluster_delete.md)	 - Delete KubeMQ cluster command
* [kubemqctl cluster describe](kubemqctl_cluster_describe.md)	 - Describe KubeMQ cluster command
* [kubemqctl cluster events](kubemqctl_cluster_events.md)	 - Show KubeMQ cluster events command
* [kubemqctl cluster get](kubemqctl_cluster_get.md)	 - Get information of KubeMQ of clusters command
* [kubemqctl cluster logs](kubemqctl_cluster_logs.md)	 - Stream logs of KubeMQ cluster pods command
* [kubemqctl cluster proxy](kubemqctl_cluster_proxy.md)	 - Proxy KubeMQ cluster connection to localhost command
* [kubemqctl cluster scale](kubemqctl_cluster_scale.md)	 - Scale KubeMQ cluster replicas command


