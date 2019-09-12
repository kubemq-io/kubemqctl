package config

var ObservabilityConfig = &EntryGroup{
	Name:    "Observability parameters",
	Entries: []Entry{
		//ObservabilityMetricsDisable,
		//ObservabilityTracingSample,
	},
	SubGroups: []*EntryGroup{
		ObservabilityPrometheusConfig,
		ObservabilityJeagerConfig,
		ObservabilityZipkinConfig,
		ObservabilityHoneycombConfig,
		ObservabilityStackDriversConfig,
		ObservabilityAWSsConfig,
		ObservabilityDatadogConfig,
	},
	Result: nil,
}

var ObservabilityMetricsDisable = &EnvEntry{
	VarName:  "METRICS_DISABLE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(Observability) Set disable observability metrics exporting:",
		Options:    []string{"false", "true"},
		Validators: nil,
		Default:    "false",
		Help:       "Set disable observability metrics exporting",
	},
}

var ObservabilityTracingSample = &EnvEntry{
	VarName:  "METRICS_TRACING_SAMPLE",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Observability) Set tracing sample probability as a percentage, i.e 0.1 =10%:",
		Validators: []Validator{IsUFloat()},
		Default:    "0.1",
		Help:       "Set tracing sample probability as a percentage",
	},
}

var ObservabilityPrometheusConfig = &EntryGroup{
	Name: "Prometheus parameters",
	Entries: []Entry{
		ObservabilityPrometheusEnable,
		ObservabilityPrometheusPath,
	},
	SubGroups: nil,
	Result:    nil,
}

var ObservabilityPrometheusEnable = &EnvEntry{
	VarName:  "METRICS_PROMETHEUS_ENABLE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(Observability-Prometheus) Enable Prometheus exporting:",
		Options:    []string{"false", "true"},
		Validators: []Validator{IsRequired()},
		Default:    "true",
		Help:       "Enable/Disable Prometheus exporting",
	},
}

var ObservabilityPrometheusPath = &EnvEntry{
	VarName:  "METRICS_PROMETHEUS_PATH",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Observability-Prometheus) Set Prometheus scraping end point (on the KubeMQ service API address):",
		Validators: []Validator{IsRequired()},
		Default:    "/metrics",
		Help:       "Set Prometheus scraping end point (on the KubeMQ service API address)",
	},
}

var ObservabilityJeagerConfig = &EntryGroup{
	Name: "Jeager parameters",
	Entries: []Entry{
		ObservabilityJeagersEnable,
		ObservabilityJeagerCollectorAddress,
		ObservabilityJeagerAGENTAddress,
	},
	SubGroups: nil,
	Result:    nil,
}

var ObservabilityJeagersEnable = &EnvEntry{
	VarName:  "METRICS_JEAGER_ENABLE",
	VarValue: "0",
	Prompt: &Selection{
		Message:    "(Observability-Jeager) Enable Jeager exporting:",
		Options:    []string{"false", "true"},
		Validators: []Validator{IsRequired()},
		Default:    "false",
		Help:       "Enable Jeager exporting",
	},
}

var ObservabilityJeagerCollectorAddress = &EnvEntry{
	VarName:  "METRICS_JEAGER_COLLECTOR_ADDRESS",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Observability-Jeager) Set Jeager collector address (host:port):",
		Validators: []Validator{},
		Default:    "",
		Help:       "Sets Jeager collector address",
	},
}

var ObservabilityJeagerAGENTAddress = &EnvEntry{
	VarName:  "METRICS_JEAGER_AGENT_ADDRESS",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Observability-Jeager) Set Jeager agent address (host:port):",
		Validators: []Validator{},
		Default:    "",
		Help:       "Sets Jeager agent address",
	},
}

var ObservabilityZipkinConfig = &EntryGroup{
	Name: "Zipkin parameters",
	Entries: []Entry{
		ObservabilityZipkinEnable,
		ObservabilityZipkinAReporterAddress,
	},
	SubGroups: nil,
	Result:    nil,
}

var ObservabilityZipkinEnable = &EnvEntry{
	VarName:  "METRICS_ZIPKIN_ENABLE",
	VarValue: "0",
	Prompt: &Selection{
		Message:    "(Observability-Zipkin) Enable Zipkin exporting (host:port):",
		Options:    []string{"false", "true"},
		Validators: []Validator{IsRequired()},
		Default:    "false",
		Help:       "Enable Zipkin exporting",
	},
}

var ObservabilityZipkinAReporterAddress = &EnvEntry{
	VarName:  "METRICS_ZIPKEIN_REPORTER_ADDRESS",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Observability-Zipkin) Set Zipkin agent address:",
		Validators: []Validator{},
		Default:    "",
		Help:       "Sets Zipkin agent address",
	},
}

var ObservabilityHoneycombConfig = &EntryGroup{
	Name: "Honeycomb parameters",
	Entries: []Entry{
		ObservabilityHoneycombsEnable,
		ObservabilityHoneycombKey,
		ObservabilityHoneycombDataset,
	},
	SubGroups: nil,
	Result:    nil,
}

var ObservabilityHoneycombsEnable = &EnvEntry{
	VarName:  "METRICS_HONEYCOMB_ENABLE",
	VarValue: "0",
	Prompt: &Selection{
		Message:    "(Observability-Honeycomb) Enable Honeycomb exporting:",
		Options:    []string{"false", "true"},
		Validators: []Validator{IsRequired()},
		Default:    "false",
		Help:       "Enable Honeycomb exporting",
	},
}

var ObservabilityHoneycombKey = &EnvEntry{
	VarName:  "METRICS_HONEYCOMB_KEY",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Observability-Honeycomb) Set Honeycomb's key:",
		Validators: []Validator{},
		Default:    "",
		Help:       "Set Honeycomb's key",
	},
}

var ObservabilityHoneycombDataset = &EnvEntry{
	VarName:  "METRICS_HONEYCOMB_DATASET",
	VarValue: "0",
	Prompt: &Input{
		Message:    "(Observability-Honeycomb) Set Honeycomb's dataset:",
		Validators: []Validator{},
		Default:    "",
		Help:       "Set Honeycomb's dataset",
	},
}

var ObservabilityStackDriversConfig = &EntryGroup{
	Name: "Google's StackDriver parameters",
	Entries: []Entry{
		ObservabilityStackDriversEnable,
		ObservabilityStackDriverProjectID,
		ObservabilityStackDriverMonitorCreds,
		ObservabilityStackDriverTraceCreds,
	},
	SubGroups: nil,
	Result:    nil,
}

var ObservabilityStackDriversEnable = &EnvEntry{
	VarName:  "METRICS_STACKDRIVER_ENABLE",
	VarValue: "0",
	Prompt: &Selection{
		Message:    "(Observability-StackDriver) Enable StackDriver exporting:",
		Options:    []string{"false", "true"},
		Validators: []Validator{IsRequired()},
		Default:    "false",
		Help:       "Enable StackDriver exporting",
	},
}

var ObservabilityStackDriverProjectID = &EnvEntry{
	VarName:  "METRICS_STACKDRIVER_PROJECT_ID",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Observability-StackDriver) Set StackDriver project id:",
		Validators: []Validator{},
		Default:    "",
		Help:       "Set StackDriver project id",
	},
}

var ObservabilityStackDriverMonitorCreds = &EnvEntry{
	VarName:  "METRICS_STACKDRIVER_MONITOR_CRDES",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Observability-StackDriver) Set StackDriver monitor (stats) credentials file location:",
		Validators: []Validator{},
		Default:    "./creds.pem",
		Help:       "Set StackDriver monitor(stats)credentials file location",
	},
}

var ObservabilityStackDriverTraceCreds = &EnvEntry{
	VarName:  "METRICS_STACKDRIVER_TRACE_CREDS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Observability-StackDriver) Set StackDriver traces credentials file location:",
		Validators: []Validator{},
		Default:    "./creds.pem",
		Help:       "Set StackDriver traces credentials file location",
	},
}

var ObservabilityAWSsConfig = &EntryGroup{
	Name: "AWS X-Ray parameters",
	Entries: []Entry{
		ObservabilityAWSsEnable,
		ObservabilityAWSAccessKeyID,
		ObservabilityAWSSecretAccessKey,
		ObservabilityAWSDefaultRegion,
	},
	SubGroups: nil,
	Result:    nil,
}

var ObservabilityAWSsEnable = &EnvEntry{
	VarName:  "METRICS_AWS_ENABLE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(Observability-AWS) Enable AWS X-Ray exporting:",
		Options:    []string{"false", "true"},
		Validators: []Validator{IsRequired()},
		Default:    "false",
		Help:       "Enable AWS X-Ray exporting",
	},
}

var ObservabilityAWSAccessKeyID = &EnvEntry{
	VarName:  "METRICS_AWS_ACCESS_KEY_ID",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Observability-AWS) Set AWS access key id environment variable:",
		Validators: []Validator{},
		Default:    "",
		Help:       "Set AWS access key id environment variable",
	},
}

var ObservabilityAWSSecretAccessKey = &EnvEntry{
	VarName:  "METRICS_AWS_SECRET_ACCESS_KEY",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Observability-AWS) Set AWS secret access key environment variable:",
		Validators: []Validator{},
		Default:    "",
		Help:       "Set AWS secret access key environment variable",
	},
}

var ObservabilityAWSDefaultRegion = &EnvEntry{
	VarName:  "METRICS_AWS_DEFAULT_REGION",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Observability-AWS) Set AWS default region environment variable:",
		Validators: []Validator{},
		Default:    "",
		Help:       "Set AWS default region environment variable",
	},
}

var ObservabilityDatadogConfig = &EntryGroup{
	Name: "Datadog parameters",
	Entries: []Entry{
		ObservabilityDatadogsEnable,
		ObservabilityDatadogTraceAddress,
		ObservabilityDatadogStatsAddress,
	},
	SubGroups: nil,
	Result:    nil,
}

var ObservabilityDatadogsEnable = &EnvEntry{
	VarName:  "METRICS_DATADOG_ENABLE",
	VarValue: "",
	Prompt: &Selection{
		Message:    "(Observability-Datadog) Enable Datadog exporting:",
		Options:    []string{"false", "true"},
		Validators: []Validator{IsRequired()},
		Default:    "false",
		Help:       "Enable Datadog exporting",
	},
}

var ObservabilityDatadogTraceAddress = &EnvEntry{
	VarName:  "METRICS_DATADOG_TRACE_ADDRESS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Observability-Datadog) Set Datadog's trace address (host:port):",
		Validators: []Validator{},
		Default:    "",
		Help:       "Set Datadog's trace address",
	},
}

var ObservabilityDatadogStatsAddress = &EnvEntry{
	VarName:  "METRICS_DATADOG_STATS_ADDRESS",
	VarValue: "",
	Prompt: &Input{
		Message:    "(Observability-Datadog) Set Datadog's stats address (host:port):",
		Validators: []Validator{},
		Default:    "",
		Help:       "Sets Datadog's stats address",
	},
}
