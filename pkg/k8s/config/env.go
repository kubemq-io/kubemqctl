package config

var CreateEnvGroup = &EntryGroup{
	Name:    "Environment variables",
	Entries: nil,
	SubGroups: []*EntryGroup{
		LicenseConfig,
		PersistenceConfig,
		QueuesConfig,
		GrpcConfig,
		RESTConfig,
		LoggingConfig,
		ObservabilityConfig,
	},
	Result: nil,
}
