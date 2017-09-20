package zbmsgpack

type Topic struct {
	Name       string `msgpack:"name"`
	State      string `msgpack:"state"`
	Partitions int    `msgpack:"partitions"`
}
