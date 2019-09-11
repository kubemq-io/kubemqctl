package config

var AuthenticationGRPCConfig = &EntryGroup{
	Name: "",
	Entries: []Entry{
		&NoPromptEntry{
			VarName:  "GRPC_SECURITY_TLS_MODE",
			VarValue: "tls",
		},
		&NoPromptEntry{
			VarName:  "GRPC_SECURITY_CERT_FILE",
			VarValue: "./tls-grpc.cert",
		},
		&NoPromptEntry{
			VarName:  "GRPC_SECURITY_KEY_FILE",
			VarValue: "./tls-grpc.key",
		},
		&SecretEntry{
			Name:        "tls-grpc-cert",
			Description: "gRPC TLS Cert file",
			ClusterName: CreateBasicOptions.Name,
			EnvVarName:  "GRPC_SECURITY_CERT_FILE",
			FilePath:    "./tls-grpc.cert",
			FileName:    "tls-grpc.cert",
			SecretType:  "tls.cert",
		},
		&SecretEntry{
			Name:        "tls-grpc-key",
			Description: "gRPC TLS Key file",
			ClusterName: CreateBasicOptions.Name,
			EnvVarName:  "GRPC_SECURITY_KEY_FILE",
			FilePath:    "./tls-grpc.key",
			FileName:    "tls-grpc.key",
			SecretType:  "tls.key",
		},
	},
	SubGroups: nil,
	Result:    nil,
}

var AuthenticationRESTConfig = &EntryGroup{
	Name: "",
	Entries: []Entry{
		&NoPromptEntry{
			VarName:  "REST_SECURITY_TLS_MODE",
			VarValue: "tls",
		},
		&NoPromptEntry{
			VarName:  "REST_SECURITY_CERT_FILE",
			VarValue: "./tls-rest.cert",
		},
		&NoPromptEntry{
			VarName:  "REST_SECURITY_KEY_FILE",
			VarValue: "./tls-rest.key",
		},
		&SecretEntry{
			Name:        "tls-rest-cert",
			Description: "REST TLS Cert file",
			ClusterName: CreateBasicOptions.Name,
			EnvVarName:  "REST_SECURITY_CERT_FILE",
			FilePath:    "./tls-rest.cert",
			FileName:    "tls-rest.cert",
			SecretType:  "tls.cert",
		},
		&SecretEntry{
			Name:        "tls-rest-key",
			Description: "REST TLS Key file",
			ClusterName: CreateBasicOptions.Name,
			EnvVarName:  "REST_SECURITY_KEY_FILE",
			FilePath:    "./tls-rest.key",
			FileName:    "tls-rest.key",
			SecretType:  "tls.key",
		},
	},
	SubGroups: nil,
	Result:    nil,
}
