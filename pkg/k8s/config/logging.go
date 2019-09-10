package config

var LoggingConfig = &EntryGroup{
	Name: "Logging parameters",
	Entries: []Entry{
		LoggingLogLevel,
		LoggingEnabelLogFile,
		LoggingLogFilePath,
	},
	SubGroups: nil,
	Result:    nil,
}

var LoggingLogLevel = &EnvEntry{
	VarName:  "KUBEMQ_LOG_LEVEL",
	VarValue: "2",
	Prompt: &Selection{
		Message:    "(Logging) Select KubeMQ stdout log level:",
		Options:    []string{"1", "2", "3", "4", "5"},
		Validators: nil,
		Default:    "2",
		Help:       "1-Debug, 2-Info, 3-Warn, 4-Error, 5-Fatal",
	},
}

var LoggingEnabelLogFile = &EnvEntry{
	VarName:  "LOG_FILE_ENABLE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(Logging) Enable saving logs to file",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "false",
		Help:       "Enable/Disable saving logs to file",
	},
}

var LoggingLogFilePath = &EnvEntry{
	VarName:  "LOG_FILE_PATH",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Logging-File) Set log file write path:",
		Validators: nil,
		Default:    "./log",
		Help:       "Set file write path, default: ./log",
	},
}

var LoggingEnabelLoggly = &EnvEntry{
	VarName:  "LOG_LOGGLY_ENABLE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(Logging) Enable sending logs to https://www.loggly.com/ external service",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "false",
		Help:       "Enable/Disable sending logs to https://www.loggly.com/ external service",
	},
}

var LoggingLogglyKey = &EnvEntry{
	VarName:  "LOG_LOGGLY_KEY",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Logging-Loggly) Set Loggly access key:",
		Validators: nil,
		Default:    "",
		Help:       "Set Loggly access key",
	},
}

var LoggingLogglyFlushInterval = &EnvEntry{
	VarName:  "LOG_LOGGLY_FLUSH_INTERVAL",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Logging-Loggly) Set Loggly sending logs interval in seconds:",
		Validators: nil,
		Default:    "5",
		Help:       "Set Loggly sending logs interval in seconds",
	},
}
