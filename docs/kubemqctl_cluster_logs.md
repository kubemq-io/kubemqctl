## kubemqctl cluster logs

Stream logs of Kubemq cluster pods command

### Synopsis

Stream command allows to show pods logs with powerful filtering capabilities

```
kubemqctl cluster logs [flags]
```

### Examples

```

	# Stream logs with selection of Kubemq cluster
	kubemqctl cluster logs

	# Stream logs of all pods in default namespace
	kubemqctl cluster logs .* -n default

	# Stream logs of regex base pods with logs since 10m ago
	kubemqctl cluster logs kubemq-cluster.* -s 10m

	# Stream logs of regex base pods with logs since 10m ago include the string of 'connection'
	kubemqctl cluster logs kubemq-cluster.* -s 10m -i connection

	# Stream logs of regex base pods with logs exclude the string of 'error'
	kubemqctl cluster logs kubemq-cluster.* -s 10m -e error

	# Stream logs of specific container
	kubemqctl cluster logs -c kubemq-cluster-0

```

### Options

```
  -c, --container string      Set container regex
      --disable-color         Set to disable colorized output
  -e, --exclude stringArray   Set strings to exclude
  -h, --help                  help for logs
  -i, --include stringArray   Set strings to include
  -l, --label string          Set label selector
  -n, --namespace string      Set default namespace
  -s, --since duration        Set since duration time
  -t, --tail int              Set how many lines to tail for each pod
```

### SEE ALSO

* [kubemqctl cluster](kubemqctl_cluster.md)	 - Executes Kubemq cluster management commands


