package config

var RESTInterfaceConfig = &EntryGroup{
	Name: "REST interface parameters",
	Entries: []Entry{
		RESTEnable,
		RESTSubBufSize,
		RESTBodyLimit,
		RESTReadTimeout,
		RESTWriteTimeout,
	},
	Result: nil,
}

var RESTEnable = &EnvEntry{
	VarName:  "REST_ENABLE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(REST) Disable REST interface:",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "false",
		Help:       "",
	},
}

var RESTSubBufSize = &EnvEntry{
	VarName:  "REST_SUB_BUFF_SIZE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST) Set the subscribe message / requests buffer size to use on the server:",
		Validators: []Validator{IsUint()},
		Default:    "100",
		Help:       "",
	},
}

var RESTBodyLimit = &EnvEntry{
	VarName:  "REST_BODY_LIMIT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST) Set request body limit, (i.e. 2M):",
		Validators: nil,
		Default:    "",
		Help:       "",
	},
}

var RESTReadTimeout = &EnvEntry{
	VarName:  "REST_READ_TIMEOUT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-Security) Set REST read timeout in seconds:",
		Validators: nil,
		Default:    "60",
		Help:       "",
	},
}
var RESTWriteTimeout = &EnvEntry{
	VarName:  "REST_WRITE_TIMEOUT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-Security) Set REST write timeout in seconds:",
		Validators: []Validator{IsUint()},
		Default:    "60",
		Help:       "Set REST write timeout in seconds",
	},
}

var RESTInterfaceCORS = &EntryGroup{
	Name: "REST CORS parameters",
	Entries: []Entry{
		RESTAllowOrigins,
		RESTAllowMethods,
		RESTAllowHeaders,
		RESTAllowCredentials,
		RESTExposeHeaders,
		RESTMaxAge,
	},
	SubGroups: nil,
	Result:    nil,
}

var RESTAllowOrigins = &EnvEntry{
	VarName:  "REST_CORS_ALLOW_ORIGINS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Defines a list of origins that may access the resource:",
		Validators: []Validator{IsRequired()},
		Default:    "{*}",
		Help:       "",
	},
}
var RESTAllowMethods = &EnvEntry{
	VarName:  "REST_CORS_ALLOW_METHODS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Set a list of origins that may access the resource:",
		Validators: []Validator{IsRequired()},
		Default:    `{"GET", "POST"}`,
		Help:       "",
	},
}
var RESTAllowHeaders = &EnvEntry{
	VarName:  "REST_CORS_ALLOW_HEADERS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Set a list of request headers that can be used when making the actual request:",
		Validators: []Validator{IsRequired()},
		Default:    "{}",
		Help:       "",
	},
}
var RESTAllowCredentials = &EnvEntry{
	VarName:  "REST_CORS_ALLOW_CREDENTIALS",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(REST-CORS) Set whether or not the response to the request can be exposed when the credentials flag is true:",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "false",
		Help:       "",
	},
}
var RESTExposeHeaders = &EnvEntry{
	VarName:  "REST_CORS_EXPOSE_HEADERS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Set a whitelist headers that clients are allowed to access:",
		Validators: []Validator{IsRequired()},
		Default:    "{}",
		Help:       "",
	},
}
var RESTMaxAge = &EnvEntry{
	VarName:  "REST_CORS_MAX_AGE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Sets how long (in seconds) the results of a pre-flight request can be cached:",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "",
	},
}
