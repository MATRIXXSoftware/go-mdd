package mdd

type FieldType int

const (
	Unknown FieldType = iota
	String
	Bool
	Int8
	Int16
	Int32
	Int64
	Int128
	UInt8
	UInt16
	UInt32
	UInt64
	UInt128
	Struct
)

type Containers struct {
	Containers []Container
}

type Container struct {
	Header Header
	Fields []Field
}

type Header struct {
	Version       int
	TotalField    int
	Depth         int
	Key           int
	SchemaVersion int
	ExtVersion    int
}

type Field struct {
	Data  []byte
	Type  FieldType
	Value interface{}
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
