package config

var PersistenceConfig = &EntryGroup{
	Name: "Persistence parameters",
	Entries: []*Entry{
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

var PersistenceCleanStore = &Entry{
	VarName:  "STORE_CLEAN",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(Persistence) Set clean persistence folder on start:",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "false",
		Help:       "If true, KubeMQ will clean all the files in the store on boot",
	},
}

var PersistenceMaxQueue = &Entry{
	VarName:  "STORE_MAX_QUEUES",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Persistence) Set max number of persistent channels/queues (0 - unlimited):",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "Set max number of persistent channels/queues, 0 = unlimited",
	},
}

var PersistenceMaxSubscribers = &Entry{
	VarName:  "STORE_MAX_SUBSCRIBERS",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Persistence) Set max number of subscribers per channel/queue (0 - unlimited):",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "Set max number of subscribers per channel/queue, 0 = unlimited",
	},
}

var PersistenceMaxMessages = &Entry{
	VarName:  "STORE_MAX_MESSAGES",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Persistence) Set max number of stored messages per channel/queue (0 - unlimited):",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "Set max number of stored messages per channel/queue, 0 = unlimited",
	},
}
var PersistenceMaxSize = &Entry{
	VarName:  "STORE_MAX_SIZE",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Persistence) Set KubeMQ max size per channel/queue in bytes (0 - unlimited):",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "Set max size in bytes per channel/queue, 0 = unlimited ",
	},
}
var PersistenceMaxRetention = &Entry{
	VarName:  "STORE_MAX_RETENTION",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Persistence) Set store time for each message per channel/queue in minutes (0 - infinite):",
		Validators: []Validator{IsUint()},
		Default:    "1440",
		Help:       "Set store time in minutes for each message per channel/queue, 0 = infinite",
	},
}
var PersistenceMaxPurge = &Entry{
	VarName:  "STORE_MAX_INACTIVITY_PURGE",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Persistence) Set delete channel/queue due to inactivity time in minutes (0 - infinite):",
		Validators: []Validator{IsUint()},
		Default:    "1440",
		Help:       "Set delete channel/queue due to inactivity time in minutes, 0 = no purging",
	},
}
