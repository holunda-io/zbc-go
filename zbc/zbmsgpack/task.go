package zbmsgpack

type Task struct {
	State       string                 `yaml:"state" zbmsgpack:"state"`
	Headers     map[string]interface{} `yaml:"headers" zbmsgpack:"headers"`
	Retries     int                    `yaml:"retries" zbmsgpack:"retries"`
	Type        string                 `yaml:"type" zbmsgpack:"type"`
	Payload     []uint8                `yaml:"-" zbmsgpack:"payload"`
	PayloadJSON map[string]interface{} `yaml:"payload" zbmsgpack:"-" json:"-"`
}


func NewTask(typeName string) *Task {
	return &Task{
		State: "CREATE",
		Type: typeName,
	}
}