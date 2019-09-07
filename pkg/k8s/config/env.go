package config

var EnvConfig = &EntryGroup{
	Name: "Environment variables",
	Entries: []*Entry{
		DefaultToken,
	},
	SubGroups: []*EntryGroup{
		LicenseConfig,
		LoggingConfig,
		PersistenceConfig,
	},
	Result: nil,
}
