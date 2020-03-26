## kubemqctl create cluster

Create a Kubemq cluster command

### Synopsis

Create command allows to deploy a Kubemq cluster with configuration options

```
kubemqctl create cluster [flags]
```

### Examples

```

	# Create default Kubemq cluster
	kubemqctl create cluster
	
	# Create Kubemq cluster with options - get all flags
	kubemqctl create cluster --help

```

### Options

```
      --api-disabled                               disable Api interface
      --api-expose string                          set api port service type (ClusterIP,NodePort,LoadBalancer) (default "ClusterIP")
      --api-node-port int32                        set api node port value
      --api-port int32                             set api port value (default 8080)
      --authentication-enabled                     enable authentication configuration
      --authentication-public-key-data string      set authentication public key data
      --authentication-public-key-file string      set authentication public key filename
      --authentication-public-key-type string      set authentication public key type
      --authorization-auto-reload int32            set authorization auto policy loading time interval in minutes
      --authorization-enabled                      enable authorization configuration
      --authorization-policy-data string           set authorization policy data
      --authorization-policy-file string           set authorization policy filename
      --authorization-url string                   set authorization policy loading url
  -c, --config-file string                         set kubemq config file
      --dry-run                                    generate cluster configuration without execute
      --gateway-ca-data string                     set tls ca certificate data for remote gateway
      --gateway-ca-file string                     set tls ca certificate filename for remote gateway
      --gateway-cert-data string                   set tls certificate data for remote gateway
      --gateway-cert-file string                   set tls certificate filename for remote gateway
      --gateway-enabled                            enable gateway configuration
      --gateway-key-data string                    set tls key data for remote gateway
      --gateway-key-file string                    set tls key filename for remote gateway 
      --gateway-port int32                         set gateway listen port value (default 7000)
      --gateway-remotes stringArray                set tls certificate data for remote gateway
      --grpc-body-limit int32                      set Max size of payload in bytes 
      --grpc-buffer-size int32                     set subscribe message / requests buffer size to use on server 
      --grpc-disabled                              disable grpc interface
      --grpc-expose string                         set grpc port service type (ClusterIP,NodePort,LoadBalancer) (default "ClusterIP")
      --grpc-node-port int32                       set grpc node port value
      --grpc-port int32                            set grpc port value (default 50000)
      --health-enabled                             enable resources configuration
      --health-failure-threshold int32             set health prob failure threshold (default 6)
      --health-initial-delay int32                 set health prob initial delay seconds  (default 5)
      --health-period-seconds int32                set health prob period seconds  (default 10)
      --health-success-threshold int32             set health prob success threshold (default 1)
      --health-timout-seconds int32                set health prob timeout seconds  (default 5)
  -h, --help                                       help for cluster
      --image string                               set image registry/repository:tag (default "docker.io/kubemq/kubemq:latest")
      --image-pull-policy string                   set image pull policy (default "Always")
      --license-data string                        set license data
      --license-filename string                    set license filename
  -t, --license-token string                       set license token
      --log-data int32                             set log level (default 2)
      --log-file string                            set log filename
      --name string                                set kubemq cluster name (default "kubemq-cluster")
  -n, --namespace string                           set kubemq cluster namespace (default "kubemq")
      --node-selectors-keys stringToString         set statefulset node selectors key-value (map) (default [])
      --notification-enabled                       set notification enable
      --notification-log                           set log notification to std-out
      --notification-prefix string                 set notification channel prefix
      --queue-default-visibility-seconds int32     set default time of hold received message before returning to queue (default 60)
      --queue-default-wait-timeout-seconds int32   set default time to wait for a message in a queue (default 1)
      --queue-max-delay-seconds int32              set max delay seconds allowed for message (default 43200)
      --queue-max-expiration-seconds int32         set max expiration allowed for message (default 43200)
      --queue-max-receive-messages-request int32   set max of sending / receiving batch of queue message  (default 1024)
      --queue-max-requeues int32                   set max retires to receive message before discard (default 1024)
      --queue-max-visibility-seconds int32         set max time of hold received message before returning to queue (default 43200)
      --queue-max-wait-timeout-seconds int32       set max wait timeout allowed for message (default 3600)
  -r, --replicas int32                             set replicas (default 3)
      --resources-enabled                          enable resources configuration
      --resources-limits-key-cpu string            set resources limits cpu  (default "1000m")
      --resources-limits-key-memory string         set resources limits memory (default "512Mi")
      --resources-requests-key-cpu string          set resources requests cpu (default "100m")
      --resources-requests-memory string           set resources request memory (default "256Mi")
      --rest-body-limit int32                      set Max size of payload in bytes 
      --rest-buffer-size int32                     set subscribe message / requests buffer size to use on server 
      --rest-disabled                              disable rest interface
      --rest-expose string                         set rest port service type (ClusterIP,NodePort,LoadBalancer) (default "ClusterIP")
      --rest-node-port int32                       set rest node port value
      --rest-port int32                            set rest port value (default 9090)
      --routing-auto-reload int32                  set routing auto loading time interval in minutes
      --routing-data string                        set routing data
      --routing-filename string                    set routing filename
      --routing-url string                         set routing loading url
      --store-clean                                set clear persistence data on start-up    
      --store-max-channel-size int32               Set limit size of channel in bytes 
      --store-max-channels int32                   set limit number of persistence channels
      --store-max-messages int32                   set limit of messages per channel
      --store-max-subscribers int32                set limit of subscribers per channel
      --store-messages-retention-minutes int32     set message retention time in minutes (default 1440)
      --store-path string                          set persistence file path (default "./store")
      --store-purge-inactive-minutes int32         set time in minutes of channel inactivity to delete (default 1440)
      --tls-ca-data string                         set tls ca certificate data
      --tls-ca-file string                         set tls ca certificate filename
      --tls-cert-data string                       set tls certificate data
      --tls-cert-file string                       set tls certificate filename
      --tls-enabled                                enable tls tls configuration
      --tls-key-data string                        set tls key data
      --tls-key-file string                        set tls key filename
  -v, --volume-size string                         set persisted volume size
```

### Options inherited from parent commands

```
      --config string   set kubemqctl configuration file (default "./.kubemqctl.yaml")
```

### SEE ALSO

* [kubemqctl create](kubemqctl_create.md)	 - Executes Kubemq create commands

###### Auto generated by spf13/cobra on 26-Mar-2020
