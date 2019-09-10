package config

var GrpcConfig = &EntryGroup{
	Name: "gRPC interface parameters",
	Entries: []Entry{
		GrpcEnable,
		GrpcSubBufSize,
		GrpcBodyLimit,
		GrpcTLSMode,
		GrpcCertFile,
		GrpcKeyFile,
	},
	SubGroups: nil,
	Result:    nil,
}

var GrpcEnable = &EnvEntry{
	VarName:  "GRPC_ENABLE",
	VarValue: "0",
	Prompt: &Selection{
		Message:    "(gRPC) Enable gRPC interface:",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "true",
		Help:       "Enable/Disable the gRPC interface",
	},
}

var GrpcSubBufSize = &EnvEntry{
	VarName:  "GRPC_SUB_BUFF_SIZE",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(gRPC) Set the subscribe message / requests buffer size to use on the serve:",
		Validators: []Validator{IsUint()},
		Default:    "100",
		Help:       "Set the subscribe message / requests buffer size to use on the serve",
	},
}

var GrpcBodyLimit = &EnvEntry{
	VarName:  "GRPC_BODY_LIMIT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(gRPC) Set request body limit in bytes:",
		Validators: []Validator{IsUint()},
		Default:    "4194304",
		Help:       "Set request body limit in bytes",
	},
}

var GrpcTLSMode = &EnvEntry{
	VarName:  "GRPC_SECURITY_TLS_MODE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(gRPC-Security) Set Security mode:",
		Options:    []string{"none", "tls"},
		Validators: nil,
		Default:    "none",
		Help:       "Set Security mode",
	},
}

var GrpcCertFile = &EnvEntry{
	VarName:  "GRPC_SECURITY_CERT_FILE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(gRPC-Security) Set CERT file name and location:",
		Validators: nil,
		Default:    "./cert.pem",
		Help:       "Set CERT file name and location",
	},
}

var GrpcKeyFile = &EnvEntry{
	VarName:  "GRPC_SECURITY_KEY_FILE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(gRPC-Security) Set Key file name and location:",
		Validators: nil,
		Default:    "./key.pem",
		Help:       "Set Key file name and location",
	},
}
