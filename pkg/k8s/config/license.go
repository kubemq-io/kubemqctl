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
		Message:    "Set KubeMQ license data:",
		Validators: nil,
		Default:    "",
		Help:       "Sets KubeMQ license data",
	},
}

var LicenseProxy = &Entry{
	VarName:  "KUBEMQ_PROXY",
	VarValue: "",
	Prompt: &Input{
		Message:    "Set Proxy server address url access:",
		Validators: []Validator{IsValidHostPort()},
		Default:    "",
		Help:       "Sets Proxy server address url access (in case license validation failure)",
	},
}
