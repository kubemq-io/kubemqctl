## kubemqctl cluster

Executes Kubemq cluster management commands

### Synopsis

Executes Kubemq cluster management commands

```
kubemqctl cluster [flags]
```

### Examples

```

	# Execute create Kubemq cluster command
	kubemqctl cluster create

	# Execute delete Kubemq cluster command
	kubemqctl cluster delete

	# Execute describe Kubemq cluster command
	kubemqctl cluster describe

	# Execute apply Kubemq cluster command
	kubemqctl cluster apply

	# Execute show Kubemq cluster logs command
	kubemqctl cluster logs

	# Execute scale Kubemq cluster command
	kubemqctl cluster scale

	# Execute list of Kubemq clusters command
	kubemqctl cluster list

	# Execute proxy ports of Kubemq cluster command
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
* [kubemqctl cluster apply](kubemqctl_cluster_apply.md)	 - Apply a Kubemq cluster command
* [kubemqctl cluster context](kubemqctl_cluster_context.md)	 - Select kubernetes cluster context command
* [kubemqctl cluster create](kubemqctl_cluster_create.md)	 - Create a Kubemq cluster command
* [kubemqctl cluster dashboard](kubemqctl_cluster_dashboard.md)	 - Dashboard command allows to start a web view of Kubemq cluster dashboard
* [kubemqctl cluster delete](kubemqctl_cluster_delete.md)	 - Delete Kubemq cluster command
* [kubemqctl cluster describe](kubemqctl_cluster_describe.md)	 - Describe Kubemq cluster command
* [kubemqctl cluster events](kubemqctl_cluster_events.md)	 - Show Kubemq cluster events command
* [kubemqctl cluster get](kubemqctl_cluster_get.md)	 - Get information of Kubemq of clusters command
* [kubemqctl cluster logs](kubemqctl_cluster_logs.md)	 - Stream logs of Kubemq cluster pods command
* [kubemqctl cluster proxy](kubemqctl_cluster_proxy.md)	 - Proxy Kubemq cluster connection to localhost command
* [kubemqctl cluster scale](kubemqctl_cluster_scale.md)	 - Scale Kubemq cluster replicas command


