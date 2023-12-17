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

func unicode(value bool) string {
	if value {
		return "✓"
	}
	return "✗"
}
