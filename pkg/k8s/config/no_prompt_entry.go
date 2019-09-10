package config

type NoPromptEntry struct {
	VarName  string
	VarValue string
}

func (e *NoPromptEntry) Execute() error {
	return nil
}

func (e *NoPromptEntry) EnvVar() *EnvVar {
	return &EnvVar{
		Name:  e.VarName,
		Value: e.VarValue,
	}
}

func (e *NoPromptEntry) Volume() *Volume {
	return nil
}
func (e *NoPromptEntry) ConfigMap() *ConfigMap {
	return nil
}

func (e *NoPromptEntry) Secret() *Secret {
	return nil
}
