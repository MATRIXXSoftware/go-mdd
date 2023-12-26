package mdd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/matrixxsoftware/go-mdd/mdd/field"
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
	Data        []byte
	Type        field.Type
	IsMulti     bool
	IsContainer bool
	IsNull      bool
	Value       interface{}
	Codec       Codec
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

// Dump to string

func (c *Containers) Dump() string {
	var sb strings.Builder
	for _, container := range c.Containers {
		sb.WriteString(container.Dump())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (c *Container) Dump() string {
	var sb strings.Builder
	sb.WriteString(c.Header.Dump())

	sb.WriteString(fmt.Sprintf(" %5s  %-10s %8s %8s   %-30s\n", "index", "type", "multi", "struct", "data"))
	for i, field := range c.Fields {
		sb.WriteString(fmt.Sprintf(" %5d  %-10s %8s %8s   %-30s\n", i, field.Type.String(), unicode(field.IsMulti), unicode(field.IsContainer), field.String()))
	}
	return sb.String()
}

func (h *Header) Dump() string {
	return fmt.Sprintf("%s (%d)  %d/%d\n", "Unknown", h.Key, h.SchemaVersion, h.ExtVersion)
}

func (f *Field) String() string {
	return string(f.Data)
}

func (f *Field) GetValue() (interface{}, error) {
	if f.IsNull {
		return nil, nil
	}
	if f.Value == nil {
		v, err := f.Codec.DecodeField(f)
		if err != nil {
			return nil, err
		}
		f.Value = v
	}
	return f.Value, nil
}

func unicode(value bool) string {
	if value {
		return "✓"
	}
	return "✗"
}

// Global functions to initialize Fields

func NewBasicField[T field.BasicTypes](value T) *Field {
	var fieldType field.Type

	valueType := reflect.TypeOf(value)
	switch valueType.Kind() {
	case reflect.String:
		fieldType = field.String
	case reflect.Bool:
		fieldType = field.Bool
	case reflect.Int8:
		fieldType = field.Int8
	case reflect.Int16:
		fieldType = field.Int16
	case reflect.Int32:
		fieldType = field.Int32
	case reflect.Int64:
		fieldType = field.Int64
	case reflect.Uint8:
		fieldType = field.UInt8
	case reflect.Uint16:
		fieldType = field.UInt16
	case reflect.Uint32:
		fieldType = field.UInt32
	case reflect.Uint64:
		fieldType = field.UInt64

	default:
		fieldType = field.Unknown
	}

	return &Field{
		Type:  fieldType,
		Value: value,
	}
}

func NewStructField(codec Codec, value *Containers) *Field {
	return &Field{
		Type:  field.Struct,
		Value: value,
		Codec: codec,
	}
}

func NewNullField(fieldType field.Type) *Field {
	return &Field{
		Type:   fieldType,
		IsNull: true,
	}
}
