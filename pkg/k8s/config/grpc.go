package config

var GrpcInterfaceConfig = &EntryGroup{
	Name: "gRPC interface parameters",
	Entries: []Entry{
		GrpcEnable,
		GrpcSubBufSize,
		GrpcBodyLimit,
	},
	SubGroups: nil,
	Result:    nil,
}

var GrpcEnable = &EnvEntry{
	VarName:  "GRPC_ENABLE",
	VarValue: "0",
	Prompt: &Selection{
		Message:    "(gRPC) Disable gRPC interface:",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "false",
		Help:       "",
	},
}

var GrpcSubBufSize = &EnvEntry{
	VarName:  "GRPC_SUB_BUFF_SIZE",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(gRPC) Set the subscribe message / requests buffer size to use on the serve:",
		Validators: []Validator{IsUint()},
		Default:    "100",
		Help:       "",
	},
}

var GrpcBodyLimit = &EnvEntry{
	VarName:  "GRPC_BODY_LIMIT",
	VarValue: "",
	Prompt: &Input{
		Message:    "(gRPC) Set request body limit in bytes:",
		Validators: []Validator{IsUint()},
		Default:    "4194304",
		Help:       "",
	},
}
