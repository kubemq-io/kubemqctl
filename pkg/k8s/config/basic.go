package config

var BasicClusterName = &EntryGroup{
	Name: "Basic Cluster Name",
	Entries: []Entry{&SpecEntry{
		Value: &CreateBasicOptions.Name,
		Prompt: &Input{
			Message:    "(Basic) Set cluster name:",
			Validators: []Validator{IsRequired()},
			Default:    "kubemq-cluster",
			Help:       "",
		},
	},
	},
	SubGroups: nil,
	Result:    nil,
}

var BasicClusterNamespace = &EntryGroup{
	Name: "Basic Cluster Namespace ",
	Entries: []Entry{&SpecEntry{
		Value: &CreateBasicOptions.Namespace,
		Prompt: &Input{
			Message:    "(Basic) Set cluster namespace:",
			Validators: []Validator{IsRequired()},
			Default:    "kubemq",
			Help:       "",
		},
	},
	},
	SubGroups: nil,
	Result:    nil,
}

var BasicClusterImage = &EntryGroup{
	Name: "Basic Cluster Docker Image",
	Entries: []Entry{&SpecEntry{
		Value: &CreateBasicOptions.Image,
		Prompt: &Input{
			Message:    "(Basic) Set docker image version:",
			Validators: []Validator{IsRequired()},
			Default:    "latest",
			Help:       "",
		},
	},
	},
	SubGroups: nil,
	Result:    nil,
}

var BasicClusterReplicas = &EntryGroup{
	Name: "Basic Cluster Docker Image Replicas ",
	Entries: []Entry{&SpecEntry{
		Value: &CreateBasicOptions.Replicas,
		Prompt: &Input{
			Message:    "(Basic) Set cluster node replicas:",
			Validators: []Validator{IsUint()},
			Default:    "3",
			Help:       "",
		},
	},
	},
	SubGroups: nil,
	Result:    nil,
}

var BasicClusterServiceMode = &EntryGroup{
	Name: "Basic Cluster Service Mode",
	Entries: []Entry{&SpecEntry{
		Value: &CreateBasicOptions.ServiceMode,
		Prompt: &Selection{
			Message:    "(Basic) Set cluster service mode:",
			Validators: nil,
			Options:    []string{"ClusterIP", "NodePort", "LoadBalancer"},
			Default:    "ClusterIP",
			Help:       "",
		},
	},
	},
	SubGroups: nil,
	Result:    nil,
}

var BasicClusterPVC = &EntryGroup{
	Name: "Basic Cluster PVC",
	Entries: []Entry{&SpecEntry{
		Value: &CreateBasicOptions.Vol,
		Prompt: &Input{
			Message:    "(Basic) Set cluster persistence volume claim size:",
			Validators: []Validator{IsUint()},
			Default:    "0",
			Help:       "",
		},
	},
	},
	SubGroups: nil,
	Result:    nil,
}
