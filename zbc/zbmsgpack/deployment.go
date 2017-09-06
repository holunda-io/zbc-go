package zbmsgpack

// Deployment is msg pack structure used when creating a workflow
type Deployment struct {
	State   string `yaml:"state" zbmsgpack:"state"`
	BPMNXML []byte `yaml:"bpmnXml" zbmsgpack:"bpmnXml"`
}
