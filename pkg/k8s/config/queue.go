package config

var QueuesConfig = &EntryGroup{
	Name: "Queues parameters",
	Entries: []Entry{
		QueuesMaxNumberOfMessages,
		QueuesMaxWaitTimeout,
		QueuesMaxExpirationSeconds,
		QueuesMaxDealySeconds,
		QueuesMaxReceiveCount,
		QueuesMaxVisibilitySeconds,
		QueuesDefaultVisibilitySeconds,
		QueuesDefaultWaitTimeoutSeconds,
	},
	SubGroups: nil,
	Result:    nil,
}

var QueuesMaxNumberOfMessages = &EnvEntry{
	VarName:  "QUEUE_MAX_NUMBER_OF_MESSAGE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Queues) Set max of sending / receiving batch of queue messages (0 is unlimited):",
		Validators: nil,
		Default:    "1024",
		Help:       "",
	},
}

var QueuesMaxWaitTimeout = &EnvEntry{
	VarName:  "QUEUE_MAX_WAIT_TIMEOUT_SECONDS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Queues) Set max wait time out allowed for receive message (seconds):",
		Validators: []Validator{IsUint()},
		Default:    "3600",
		Help:       "",
	},
}

var QueuesMaxExpirationSeconds = &EnvEntry{
	VarName:  "QUEUE_MAX_EXPIRATION_SECONDS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Queues) Set max expiration allowed for message (seconds):",
		Validators: []Validator{IsUint()},
		Default:    "43200",
		Help:       "Set max expiration allowed for message in seconds",
	},
}

var QueuesMaxDealySeconds = &EnvEntry{
	VarName:  "QUEUE_MAX_DELAY_SECONDS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Queues) Set max delay seconds allowed for message (seconds):",
		Validators: []Validator{IsUint()},
		Default:    "43200",
		Help:       "Set max delay seconds allowed for message",
	},
}
var QueuesMaxReceiveCount = &EnvEntry{
	VarName:  "QUEUE_MAX_RECEIVE_COUNT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Queues) Set max retires to receive message before discard:",
		Validators: []Validator{IsUint()},
		Default:    "1024",
		Help:       "",
	},
}
var QueuesMaxVisibilitySeconds = &EnvEntry{
	VarName:  "QUEUE_MAX_VISIBILITY_SECONDS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Queues) Set max time of hold received message before returning to queue (seconds):",
		Validators: []Validator{IsUint()},
		Default:    "43200",
		Help:       "",
	},
}
var QueuesDefaultVisibilitySeconds = &EnvEntry{
	VarName:  "QUEUE_DEFAULT_VISIBILITY_SECONDS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Queues) Set default time of hold received message before returning to queue (seconds):",
		Validators: []Validator{IsUint()},
		Default:    "60",
		Help:       "",
	},
}
var QueuesDefaultWaitTimeoutSeconds = &EnvEntry{
	VarName:  "QUEUE_DEFAULT_WAIT_TIMEOUT_SECONDS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Queues) Set default time to wait for a message in a queue (seconds):",
		Validators: []Validator{IsUint()},
		Default:    "1",
		Help:       "",
	},
}
