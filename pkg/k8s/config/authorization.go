package config

var AuthorizationConfig = &EntryGroup{
	Name: "",
	Entries: []Entry{
		AuthorizationACLEnable,
		AuthorizationACLData,
	},
	SubGroups: nil,
	Result:    nil,
}

var AuthorizationACLData = &ConfigMapEntry{
	Prefix:         "Authorization",
	Name:           "acl",
	Description:    "ACL (Access Control List - CSV Format)",
	ClusterName:    CreateBasicOptions.Name,
	EnvVarName:     "ACCESS_LIST",
	FilePath:       "./acl.csv",
	FileName:       "acl.csv",
	ConfigMapName:  "",
	ConfigMapValue: "",
	VolumeName:     "",
}
var AuthorizationACLEnable = &NoPromptEntry{
	VarName:  "ACCESS_LIST_ENABLE",
	VarValue: "true",
}
