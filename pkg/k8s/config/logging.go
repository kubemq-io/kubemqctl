package config

var LoggingConfig = &EntryGroup{
	Name: "Logging parameters",
	Entries: []*Entry{
		LoggingLogLevel,
		LoggingEnabelLogFile,
		LoggingLogFilePath,
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

var LoggingEnabelLogFile = &Entry{
	VarName:  "LOG_FILE_ENABLE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "Enable saving logs to file",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "false",
		Help:       "Enable/Disable saving logs to file",
	},
}

var LoggingLogFilePath = &Entry{
	VarName:  "LOG_FILE_PATH",
	VarValue: "",
	Prompt: &Input{
		Message:    "Sets log file write path:",
		Validators: nil,
		Default:    "./log",
		Help:       "Sets file write path, default: ./log",
	},
}
