package config

type EnvEntry struct {
	VarName  string
	VarValue string
	Prompt   Prompt
}

func (e *EnvEntry) Execute() error {
	if e.Prompt == nil {
		return nil
	}
	for {
		err := e.Prompt.Ask(&e.VarValue)
		if err == nil {
			return nil
		}
		if err.Error() == "interrupt" {
			return err
		}

	}
}

func (e *EnvEntry) EnvVar() *EnvVar {
	return &EnvVar{
		Name:  e.VarName,
		Value: e.VarValue,
	}
}

func (e *EnvEntry) Volume() *Volume {
	return nil
}
func (e *EnvEntry) ConfigMap() *ConfigMap {
	return nil
}
func (e *EnvEntry) Secret() *Secret {
	return nil
}
