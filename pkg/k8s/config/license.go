package config

var LicenseConfig = &EntryGroup{
	Name: "License parameters",
	Entries: []*Entry{
		LicenseData,
		LicenseProxy,
	},
	SubGroups: nil,
	Result:    nil,
}

var LicenseData = &Entry{
	VarName:  "KUBEMQ_LICENSE_DATA ",
	VarValue: "",
	Prompt: &Input{
		Message:    "(License) Set KubeMQ license data:",
		Validators: nil,
		Default:    "",
		Help:       "Set KubeMQ license data",
	},
}

var LicenseProxy = &Entry{
	VarName:  "KUBEMQ_PROXY",
	VarValue: "",
	Prompt: &Input{
		Message:    "(License) Set Proxy server address url access:",
		Validators: []Validator{IsValidHostPort()},
		Default:    "",
		Help:       "Set Proxy server address url access (in case license validation failure)",
	},
}
