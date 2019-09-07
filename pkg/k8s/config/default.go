package config

var DefaultToken = &Entry{
	VarName:  "KUBEMQ_TOKEN",
	VarValue: "",
	Prompt: &Input{
		Message:    "Set KubeMQ token key:",
		Validators: []Validator{IsRequired(), IsValidToken()},
		Default:    "",
		Help:       "Required KubeMQ Token key",
	},
}
