package zbmsgpack

type Deployment struct {
	State   string `yaml:"state" zbmsgpack:"state"`
	BPMNXML []byte `yaml:"bpmnXml" zbmsgpack:"bpmnXml"`
}