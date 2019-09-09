package config

var RESTConfig = &EntryGroup{
	Name: "REST interface parameters",
	Entries: []*Entry{
		RESTEnable,
		RESTSubBufSize,
		RESTBodyLimit,
		RESTTLSMode,
		RESTCertFile,
		RESTReadTimeout,
		RESTWriteTimeout,
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

var RESTEnable = &Entry{
	VarName:  "REST_ENABLE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(REST) Enable REST interface:",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "true",
		Help:       "Enable/Disable the REST interface",
	},
}

var RESTSubBufSize = &Entry{
	VarName:  "REST_SUB_BUFF_SIZE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST) Set the subscribe message / requests buffer size to use on the serve:",
		Validators: []Validator{IsUint()},
		Default:    "100",
		Help:       "Set the subscribe message / requests buffer size to use on the serve",
	},
}

var RESTBodyLimit = &Entry{
	VarName:  "REST_BODY_LIMIT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST) Set request body limit, (i.e. 2M):",
		Validators: nil,
		Default:    "",
		Help:       "Set request body limit, (i.e. 2M), limit can be specified as 4x or 4xB, where x is one of the multiple from K, M, G, T or P",
	},
}

var RESTTLSMode = &Entry{
	VarName:  "REST_SECURITY_TLS_MODE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(REST-Security) Set Security mode:",
		Options:    []string{"none", "tls"},
		Validators: nil,
		Default:    "none",
		Help:       "Set Security mode",
	},
}

var RESTCertFile = &Entry{
	VarName:  "REST_SECURITY_CERT_FILE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-Security) Set CERT file name and location:",
		Validators: nil,
		Default:    "./cert.pem",
		Help:       "Set CERT file name and location",
	},
}

var RESTReadTimeout = &Entry{
	VarName:  "REST_READ_TIMEOUT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-Security) Set REST read timeout in seconds:",
		Validators: nil,
		Default:    "60",
		Help:       "Set REST read timeout in seconds",
	},
}
var RESTWriteTimeout = &Entry{
	VarName:  "REST_WRITE_TIMEOUT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-Security) Set REST write timeout in seconds:",
		Validators: []Validator{IsUint()},
		Default:    "60",
		Help:       "Set REST write timeout in seconds",
	},
}

var RESTAllowOrigins = &Entry{
	VarName:  "REST_CORS_ALLOW_ORIGINS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Defines a list of origins that may access the resource:",
		Validators: []Validator{IsRequired()},
		Default:    "{*}",
		Help:       "Defines a list of origins that may access the resource",
	},
}
var RESTAllowMethods = &Entry{
	VarName:  "REST_CORS_ALLOW_METHODS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Set a list of origins that may access the resource:",
		Validators: []Validator{IsRequired()},
		Default:    `{"GET", "POST"}`,
		Help:       "Set Key file name and location",
	},
}
var RESTAllowHeaders = &Entry{
	VarName:  "REST_CORS_ALLOW_HEADERS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Set a list of request headers that can be used when making the actual request:",
		Validators: []Validator{IsRequired()},
		Default:    "{}",
		Help:       "Set a list of request headers that can be used when making the actual request",
	},
}
var RESTAllowCredentials = &Entry{
	VarName:  "REST_CORS_ALLOW_CREDENTIALS",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(REST-CORS) Set whether or not the response to the request can be exposed when the credentials flag is true:",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "false",
		Help:       "Set whether or not the response to the request can be exposed when the credentials flag is true",
	},
}
var RESTExposeHeaders = &Entry{
	VarName:  "REST_CORS_EXPOSE_HEADERS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Set a whitelist headers that clients are allowed to access:",
		Validators: []Validator{IsRequired()},
		Default:    "{}",
		Help:       "Set a whitelist headers that clients are allowed to access",
	},
}
var RESTMaxAge = &Entry{
	VarName:  "REST_CORS_MAX_AGE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Sets how long (in seconds) the results of a pre-flight request can be cached:",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "Sets how long (in seconds) the results of a pre-flight request can be cached",
	},
}