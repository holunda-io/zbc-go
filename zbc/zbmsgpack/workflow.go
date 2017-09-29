package zbmsgpack

// Workflow is msgpack structure used when creating a workflow
type Workflow struct {
	State   string `yaml:"state" msgpack:"state"`
	ResourceType string `yaml:"resourceType" msgpack:"resourceType"`
	Resource []byte `yaml:"resource" msgpack:"resource"`
}
