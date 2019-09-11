package config

type SpecEntry struct {
	Value  interface{}
	Prompt Prompt
}

func (se *SpecEntry) Execute() error {
	if se.Prompt == nil {
		return nil
	}
	for {

		err := se.Prompt.Ask(se.Value)
		if err == nil {
			return nil
		}
		if err.Error() == "interrupt" {
			return err
		}

	}
}

func (se *SpecEntry) EnvVar() *EnvVar {

	return nil
}

func (se *SpecEntry) Volume() *Volume {
	return nil
}
func (se *SpecEntry) ConfigMap() *ConfigMap {
	return nil
}
func (se *SpecEntry) Secret() *Secret {
	return nil
}
