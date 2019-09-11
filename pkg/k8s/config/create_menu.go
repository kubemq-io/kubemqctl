package config

var CreateBasicOptions = struct {
	Name        string
	Namespace   string
	Replicas    int
	Image       string
	ServiceMode string
	Vol         int
}{
	Name:        "kubemq-cluster",
	Namespace:   "default",
	Replicas:    3,
	Image:       "latest",
	ServiceMode: "ClusterIP",
	Vol:         0,
}
var CreateMenu = &Menu{
	Items: []*MenuItem{
		&MenuItem{
			Label:  "Set Basic configuration (Name, Namespace...)",
			Action: nil,
			SubMenu: &Menu{
				Items: []*MenuItem{
					&MenuItem{
						Label:  "Set cluster name (kubemq-cluster) ",
						Action: BasicClusterName,
					},
					&MenuItem{
						Label:  "Set cluster namespace (default) ",
						Action: BasicClusterNamespace,
					},
					&MenuItem{
						Label:  "Set image version (latest) ",
						Action: BasicClusterImage,
					},
					&MenuItem{
						Label:  "Set cluster node replicas (3 nodes) ",
						Action: BasicClusterReplicas,
					},
					&MenuItem{
						Label:  "Set cluster service mode (ClusterIP) ",
						Action: BasicClusterServiceMode,
					},
					&MenuItem{
						Label:  "Set PVC (Persistence Volume Claims) size (0 GiB) ",
						Action: BasicClusterPVC,
					},
				},
			},
		},
		&MenuItem{
			Label:  "Set Authentication parameters (Certs, Keys...) ",
			Action: nil,
			SubMenu: &Menu{
				Items: []*MenuItem{
					&MenuItem{
						Label:  "Set gRPC TLS Authentication",
						Action: AuthenticationGRPCConfig,
					},
					&MenuItem{
						Label:  "Set REST TLS Authentication",
						Action: AuthenticationRESTConfig,
					},
					&MenuItem{
						Label:  "Set REST CORS parameters",
						Action: RESTInterfaceCORS,
					},
				},
			},
		},
		&MenuItem{
			Label:  "Set Authorization parameters (Access Control List...) ",
			Action: AuthorizationConfig,
		},
		&MenuItem{
			Label:  "Set Messages Persistence parameters (Limits, Retentions...) ",
			Action: PersistenceConfig,
		},
		&MenuItem{
			Label:  "Set Queue Messages parameters (Limits, Retentions...) ",
			Action: QueuesConfig,
		},
		&MenuItem{
			Label:  "Set Interfaces parameters (gRPC, REST...) ",
			Action: nil,
			SubMenu: &Menu{
				Items: []*MenuItem{
					&MenuItem{
						Label:  "Set gRPC interface parameters",
						Action: GrpcInterfaceConfig,
					},
					&MenuItem{
						Label:  "Set REST interface parameters",
						Action: RESTInterfaceConfig,
					},
				},
			},
		},
		&MenuItem{
			Label:  "Set Logging parameters (Level,Exports to files...) ",
			Action: LoggingConfig,
		},
		&MenuItem{
			Label:  "Set Metrics and Tracing parameters (Prometheus, Jeager...) ",
			Action: ObservabilityConfig,
		},
		&MenuItem{
			Label:  "Set Licensing parameters (License Data, Proxy...) ",
			Action: LicenseConfig,
		},
	}}
