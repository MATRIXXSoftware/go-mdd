package mdd

import (
	"fmt"
	"strconv"
	"strings"
)

type Header struct {
	Key           string // position 0
	TotalField    int    // position 1
	Depth         int    // position 2
	SchemaVersion int    // position 3
	ExtVersion    int    // position 4
}

type Field struct {
	Value string
}

type Container struct {
	Header Header
	Fields []Field
}

func decode(data string) Container {
	fmt.Println("Decode ", data)
	var container Container

	// Extract header and fields substrings
	headerEnd := strings.Index(data, ">")
	fieldsStart := strings.Index(data, "[")

	headerStr := data[1:headerEnd]
	fieldsStr := data[fieldsStart+1 : len(data)-1]

	// Decode header
	headerParts := strings.Split(headerStr, ",")
	if len(headerParts) == 5 {
		container.Header.Key = headerParts[0]
		container.Header.TotalField, _ = strconv.Atoi(headerParts[1])
		container.Header.Depth, _ = strconv.Atoi(headerParts[2])
		container.Header.SchemaVersion, _ = strconv.Atoi(headerParts[3])
		container.Header.ExtVersion, _ = strconv.Atoi(headerParts[4])
	}

	// Decode fields
	fieldsParts := strings.Split(fieldsStr, ",")[1:] // Skipping the first part as it belongs to the header
	for _, f := range fieldsParts {
		container.Fields = append(container.Fields, Field{Value: f})
	}

	return container
}

func encode() {
	fmt.Println("Encode")
}
