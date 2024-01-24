package mdd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/matrixxsoftware/go-mdd/dictionary"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

type Containers struct {
	Containers []Container
}

type Container struct {
	Header     Header
	Fields     []Field
	Definition *dictionary.ContainerDefinition
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
	Definition  *dictionary.FieldDefinition
}

func (c *Containers) GetContainer(key int) *Container {
	for i := range c.Containers {
		container := &c.Containers[i]
		if container.Header.Key == key {
			return container
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

func (c *Container) SetField(fieldNumber int, f *Field) {
	if len(c.Fields) > fieldNumber {
		c.Fields[fieldNumber] = *f
		return
	}

	for i := len(c.Fields); i < fieldNumber; i++ {
		c.Fields = append(c.Fields, *NewNullField(field.Unknown))
	}
	c.Fields = append(c.Fields, *f)
}

func (c *Container) LoadDefinition(definition *dictionary.ContainerDefinition) {
	c.Definition = definition
	for i := range c.Fields {
		c.Fields[i].Definition = &definition.Fields[i]
		c.Fields[i].Type = definition.Fields[i].Type
	}
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

	valueType := reflect.TypeOf(value)
	fieldType := getFieldType(valueType)

	return &Field{
		Type:        fieldType,
		Value:       value,
		IsMulti:     false,
		IsContainer: false,
	}
}

func NewBasicListField[T field.BasicTypes](value []T) *Field {
	valueType := reflect.TypeOf(value).Elem()
	fieldType := getFieldType(valueType)

	return &Field{
		Type:        fieldType,
		Value:       value,
		IsMulti:     true,
		IsContainer: false,
	}
}

func NewStructField(codec Codec, value *Containers) *Field {
	return &Field{
		Type:        field.Struct,
		Value:       value,
		Codec:       codec,
		IsMulti:     false,
		IsContainer: true,
	}
}

func NewStrucListtField(codec Codec, value *Containers) *Field {
	return &Field{
		Type:        field.Struct,
		Value:       value,
		Codec:       codec,
		IsMulti:     true,
		IsContainer: true,
	}
}

func NewNullField(fieldType field.Type) *Field {
	return &Field{
		Type:   fieldType,
		IsNull: true,
	}
}

func getFieldType(t reflect.Type) field.Type {
	switch t.Kind() {
	case reflect.String:
		return field.String
	case reflect.Bool:
		return field.Bool
	case reflect.Int8:
		return field.Int8
	case reflect.Int16:
		return field.Int16
	case reflect.Int32:
		return field.Int32
	case reflect.Int64:
		return field.Int64
	case reflect.Uint8:
		return field.UInt8
	case reflect.Uint16:
		return field.UInt16
	case reflect.Uint32:
		return field.UInt32
	case reflect.Uint64:
		return field.UInt64
	default:
		return field.Unknown
	}
}
