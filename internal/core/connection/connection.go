package connection

type ConnectionTypeMeta struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	Category    string        `json:"category"`
	ConfigSpec  []ConfigField `json:"config_spec"`
}

type Source interface {
	Validate(config map[string]interface{}) error
}
