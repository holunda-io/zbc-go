package zbmsgpack

// Workflow is msgpack structure used when creating a workflow
type Workflow struct {
	State   string `yaml:"state" msgpack:"state"`
	BPMNXML []byte `yaml:"bpmnXml" msgpack:"bpmnXml"`
}
