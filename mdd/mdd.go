package mdd

type Containers struct {
	Containers []Container
}

type Container struct {
	Header Header
	Fields []Field
}

type Header struct {
	Version       int // position 0
	TotalField    int // position 1
	Depth         int // position 2
	Key           int // position 3
	SchemaVersion int // position 4
	ExtVersion    int // position 5
}

type Field struct {
	Data []byte
}

func (c *Containers) GetContainer(key int) *Container {
	for _, container := range c.Containers {
		if container.Header.Key == key {
			return &container
		}
	}
	return nil
}

func (c *Container) GetField(fieldNumber int) *Field {
	if fieldNumber >= len(c.Fields) {
		return nil
	}
	return &c.Fields[fieldNumber]
}

func (f *Field) String() string {
	return string(f.Data)
}
