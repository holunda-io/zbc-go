package zbmsgpack

type WorkflowInstance struct {
	State         string                 `yaml:"state" zbmsgpack:"state"`
	BPMNProcessID string                 `yaml:"bpmnProcessId" zbmsgpack:"bpmnProcessId"`
	Version       int                    `yaml:"version" zbmsgpack:"version"`
	Payload       []uint8                `yaml:"-" zbmsgpack:"payload"`
	PayloadJSON   map[string]interface{} `yaml:"payload" zbmsgpack:"-"`
}
