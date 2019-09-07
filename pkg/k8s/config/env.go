package config

var EnvConfig = &EntryGroup{
	Name:    "Environment variables",
	Entries: nil,
	SubGroups: []*EntryGroup{
		LicenseConfig,
		LoggingConfig,
		PersistenceConfig,
	},
	Result: nil,
}
