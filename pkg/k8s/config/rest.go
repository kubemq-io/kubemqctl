package config

var RESTConfig = &EntryGroup{
	Name: "REST interface parameters",
	Entries: []Entry{
		RESTEnable,
		RESTSubBufSize,
		RESTBodyLimit,
	},
	SubGroups: []*EntryGroup{
		RESTSec,
		RESTCORS,
	},
	Result: nil,
}

var RESTEnable = &EnvEntry{
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

var RESTSubBufSize = &EnvEntry{
	VarName:  "REST_SUB_BUFF_SIZE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST) Set the subscribe message / requests buffer size to use on the server:",
		Validators: []Validator{IsUint()},
		Default:    "100",
		Help:       "Set the subscribe message / requests buffer size to use on the serve",
	},
}

var RESTBodyLimit = &EnvEntry{
	VarName:  "REST_BODY_LIMIT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST) Set request body limit, (i.e. 2M):",
		Validators: nil,
		Default:    "",
		Help:       "Set request body limit, (i.e. 2M), limit can be specified as 4x or 4xB, where x is one of the multiple from K, M, G, T or P",
	},
}

var RESTTLSMode = &EnvEntry{
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

//var RESTCertFile = &EnvEntry{
//	EnvVarName:  "REST_SECURITY_CERT_FILE",
//	FilePath: "",
//	Prompt: &Input{
//		Message:    "(REST-Security) Set CERT file name and location:",
//		Validators: nil,
//		Default:    "./cert.pem",
//		Help:       "Set CERT file name and location",
//	},
//}

var RESTSec = &EntryGroup{
	Name: "REST Security parameters",
	Entries: []Entry{
		RESTTLSMode,
		RESTCertFile,
		RESTReadTimeout,
		RESTWriteTimeout,
	},
	SubGroups: nil,
	Result:    nil,
}
var RESTCertFile = &ConfigMapEntry{
	EnvVarName: "REST_SECURITY_CERT_FILE",
	FilePath:   "./cert.pem",
	FileName:   "cert.pem",
	VolumeName: "rest-cert",
}
var RESTReadTimeout = &EnvEntry{
	VarName:  "REST_READ_TIMEOUT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-Security) Set REST read timeout in seconds:",
		Validators: nil,
		Default:    "60",
		Help:       "Set REST read timeout in seconds",
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

var RESTCORS = &EntryGroup{
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
		Help:       "Defines a list of origins that may access the resource",
	},
}
var RESTAllowMethods = &EnvEntry{
	VarName:  "REST_CORS_ALLOW_METHODS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Set a list of origins that may access the resource:",
		Validators: []Validator{IsRequired()},
		Default:    `{"GET", "POST"}`,
		Help:       "Set Key file name and location",
	},
}
var RESTAllowHeaders = &EnvEntry{
	VarName:  "REST_CORS_ALLOW_HEADERS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Set a list of request headers that can be used when making the actual request:",
		Validators: []Validator{IsRequired()},
		Default:    "{}",
		Help:       "Set a list of request headers that can be used when making the actual request",
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
		Help:       "Set whether or not the response to the request can be exposed when the credentials flag is true",
	},
}
var RESTExposeHeaders = &EnvEntry{
	VarName:  "REST_CORS_EXPOSE_HEADERS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Set a whitelist headers that clients are allowed to access:",
		Validators: []Validator{IsRequired()},
		Default:    "{}",
		Help:       "Set a whitelist headers that clients are allowed to access",
	},
}
var RESTMaxAge = &EnvEntry{
	VarName:  "REST_CORS_MAX_AGE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(REST-CORS) Sets how long (in seconds) the results of a pre-flight request can be cached:",
		Validators: []Validator{IsUint()},
		Default:    "0",
		Help:       "Sets how long (in seconds) the results of a pre-flight request can be cached",
	},
}
