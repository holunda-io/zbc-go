package zbmsgpack

type WorkflowInstance struct {
	State         string                 `yaml:"state" msgpack:"state"`
	BPMNProcessID string                 `yaml:"bpmnProcessId" msgpack:"bpmnProcessId"`
	Version       int                    `yaml:"version" msgpack:"version"`
	Payload       []uint8                `yaml:"-" msgpack:"payload"`
	PayloadJSON   map[string]interface{} `yaml:"payload" msgpack:"-"`
}
