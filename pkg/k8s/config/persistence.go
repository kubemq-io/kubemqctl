package config

var PersistenceConfig = &EntryGroup{
	Name: "Persistence parameters",
	Entries: []Entry{
		PersistenceCleanStore,
		PersistenceMaxQueue,
		PersistenceMaxSubscribers,
		PersistenceMaxMessages,
		PersistenceMaxSize,
		PersistenceMaxRetention,
		PersistenceMaxPurge,
	},
	SubGroups: nil,
	Result:    nil,
}

var PersistenceCleanStore = &EnvEntry{
	VarName:  "STORE_CLEAN",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(Persistence) Set clean persistence folder on start (false - no clean):",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "false",
		Help:       "",
	},
}

var PersistenceMaxQueue = &EnvEntry{
	VarName:  "STORE_MAX_QUEUES",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Persistence) Set max number of persistent channels/queues (0 - unlimited):",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "",
	},
}

var PersistenceMaxSubscribers = &EnvEntry{
	VarName:  "STORE_MAX_SUBSCRIBERS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Persistence) Set max number of subscribers per channel/queue (0 - unlimited):",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "",
	},
}

var PersistenceMaxMessages = &EnvEntry{
	VarName:  "STORE_MAX_MESSAGES",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Persistence) Set max number of stored messages per channel/queue (0 - unlimited):",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "",
	},
}
var PersistenceMaxSize = &EnvEntry{
	VarName:  "STORE_MAX_SIZE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Persistence) Set KubeMQ max size per channel/queue in bytes (0 - unlimited):",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "",
	},
}
var PersistenceMaxRetention = &EnvEntry{
	VarName:  "STORE_MAX_RETENTION",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Persistence) Set store time for each message per channel/queue in minutes (0 - infinite):",
		Validators: []Validator{IsUint()},
		Default:    "1440",
		Help:       "",
	},
}
var PersistenceMaxPurge = &EnvEntry{
	VarName:  "STORE_MAX_INACTIVITY_PURGE",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Persistence) Set delete channel/queue due to inactivity time in minutes (0 - infinite):",
		Validators: []Validator{IsUint()},
		Default:    "1440",
		Help:       "",
	},
}
