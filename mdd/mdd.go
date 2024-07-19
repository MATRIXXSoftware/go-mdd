package mdd

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/matrixxsoftware/go-mdd/dictionary"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
	log "github.com/sirupsen/logrus"
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

func (c *Containers) LoadDefinition(definitions *dictionary.Dictionary) {
	for i := range c.Containers {
		container := &c.Containers[i]
		definition, found, err := definitions.Lookup(
			container.Header.Key,
			container.Header.SchemaVersion,
			container.Header.ExtVersion,
		)
		if !found {
			log.Errorf("Error lookup definition key %d schemaVersion %d extVersion %d: %v\n",
				container.Header.Key,
				container.Header.SchemaVersion,
				container.Header.ExtVersion,
				err)
		}
		container.LoadDefinition(definition)
	}

}

func (c *Containers) CastVersion(definitions *dictionary.Dictionary, schemaVersion int, extVersion int) (*Containers, error) {
	newContainers := &Containers{}
	for i := range c.Containers {
		container := &c.Containers[i]
		targetDefinition, found, err := definitions.Lookup(
			container.Header.Key,
			schemaVersion,
			extVersion,
		)
		if !found {
			return nil, fmt.Errorf("Error lookup definition key %d schemaVersion %d extVersion %d: %v",
				container.Header.Key,
				schemaVersion,
				extVersion,
				err)
		}

		newContainer, err := container.CastVersion(targetDefinition)
		if err != nil {
			return nil, err
		}

		newContainers.Containers = append(newContainers.Containers, *newContainer)
	}
	return newContainers, nil
}

func (c *Container) GetField(fieldNumber int) *Field {
	if fieldNumber >= len(c.Fields) {
		return NewNullField(field.Unknown)
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

func (c *Container) GetFieldByName(fieldName string) *Field {
	if c.Definition == nil {
		return nil
	}

	for i := range c.Fields {
		if c.Definition.Fields[i].Name == fieldName {
			return &c.Fields[i]
		}
	}
	return nil
}

func (c *Container) LoadDefinition(definition *dictionary.ContainerDefinition) {
	c.Definition = definition
	for i := range c.Fields {
		c.Fields[i].Definition = &definition.Fields[i]
		c.Fields[i].Type = definition.Fields[i].Type
	}
}

func (c *Container) CastVersion(targetDefinition *dictionary.ContainerDefinition) (*Container, error) {
	// Validation
	if targetDefinition.Key != c.Header.Key {
		return nil, fmt.Errorf("key mismatch: %d != %d", targetDefinition.Key, c.Header.Key)
	}

	if c.Definition == nil {
		return nil, fmt.Errorf("source container has no definition")
	}

	// Header
	header := Header{
		Version:       c.Header.Version,
		TotalField:    len(targetDefinition.Fields),
		Depth:         c.Header.Depth,
		Key:           c.Header.Key,
		SchemaVersion: targetDefinition.SchemaVersion,
		ExtVersion:    targetDefinition.ExtVersion,
	}

	// Fields
	fields := make([]Field, len(targetDefinition.Fields))
	for i := range targetDefinition.Fields {
		fieldName := targetDefinition.Fields[i].Name
		f := c.GetFieldByName(fieldName)
		if f != nil {
			fields[i] = *f
		} else {
			fields[i] = *NewNullField(targetDefinition.Fields[i].Type)
		}
	}

	return &Container{
		Header:     header,
		Fields:     fields,
		Definition: targetDefinition,
	}, nil
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
	if c.Definition != nil {
		// Dump with definition
		sb.WriteString(c.Header.DumpWithName(c.Definition.Name))
		sb.WriteString(fmt.Sprintf(" %5s  %-30s %-10s %8s %8s   %-30s\n",
			"index",
			"name",
			"type",
			"multi",
			"struct",
			"data"))

		for i, field := range c.Fields {
			name := c.Definition.Fields[i].Name
			value, err := field.GetValueString()
			if err != nil {
				value = field.String()
			}

			sb.WriteString(fmt.Sprintf(" %5d  %-30s %-10s %8s %8s   %-30s\n",
				i,
				name,
				field.Type.String(),
				unicode(field.IsMulti),
				unicode(field.IsContainer),
				value,
			))
		}
	} else {
		// Dump without definition
		sb.WriteString(c.Header.Dump())
		sb.WriteString(fmt.Sprintf(" %5s  %-10s %8s %8s   %-30s\n",
			"index",
			"type",
			"multi",
			"struct",
			"data"))
		for i, field := range c.Fields {
			sb.WriteString(fmt.Sprintf(" %5d  %-10s %8s %8s   %-30s\n",
				i,
				field.Type.String(),
				unicode(field.IsMulti),
				unicode(field.IsContainer),
				field.String()))
		}
	}
	return sb.String()
}

func (h *Header) DumpWithName(name string) string {
	return fmt.Sprintf("%s (%d)  %d/%d\n",
		name,
		h.Key,
		h.SchemaVersion,
		h.ExtVersion)
}

func (h *Header) Dump() string {
	return fmt.Sprintf("%s (%d)  %d/%d\n",
		"Unknown",
		h.Key,
		h.SchemaVersion,
		h.ExtVersion)
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

func (f *Field) GetValueString() (string, error) {
	v, err := f.GetValue()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", v), nil
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
