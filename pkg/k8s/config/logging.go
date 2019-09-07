package config

var LoggingConfig = &EntryGroup{
	Name: "Logging parameters",
	Entries: []*Entry{
		LoggingLogLevel,
	},
	SubGroups: nil,
	Result:    nil,
}

var LoggingLogLevel = &Entry{
	VarName:  "KUBEMQ_LOG_LEVEL",
	VarValue: "2",
	Prompt: &Selection{
		Message:    "Select KubeMQ stdout log level:",
		Options:    []string{"1", "2", "3", "4", "5"},
		Validators: nil,
		Default:    "2",
		Help:       "1-Debug, 2-Info, 3-Warn, 4-Error, 5-Fatal",
	},
}
