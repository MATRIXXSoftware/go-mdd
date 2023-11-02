package mdd

import (
	"fmt"
	"strings"
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
	Type        FieldType
	Value       interface{}
	IsMulti     bool
	IsContainer bool
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

func (ft FieldType) String() string {
	switch ft {
	case Unknown:
		return "Unknown"
	case String:
		return "String"
	case Bool:
		return "Bool"
	case Int8:
		return "Int8"
	case Int16:
		return "Int16"
	case Int32:
		return "Int32"
	case Int64:
		return "Int64"
	case Int128:
		return "Int128"
	case UInt8:
		return "UInt8"
	case UInt16:
		return "UInt16"
	case UInt32:
		return "UInt32"
	case UInt64:
		return "UInt64"
	case UInt128:
		return "UInt128"
	case Struct:
		return "Struct"
	default:
		return "Undefined"
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

func unicode(value bool) string {
	if value {
		return "✓"
	}
	return "✗"
}
