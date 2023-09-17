package mdd

import (
	"fmt"
	"strconv"
	"strings"
)

type Header struct {
	Version       int // position 0
	TotalField    int // position 1
	Depth         int // position 2
	Key           int // position 3
	SchemaVersion int // position 4
	ExtVersion    int // position 5
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
	if len(headerParts) == 6 {
		container.Header.Version, _ = strconv.Atoi(headerParts[0])
		container.Header.TotalField, _ = strconv.Atoi(headerParts[1])
		container.Header.Depth, _ = strconv.Atoi(headerParts[2])
		container.Header.Key, _ = strconv.Atoi(headerParts[3])
		container.Header.SchemaVersion, _ = strconv.Atoi(headerParts[4])
		container.Header.ExtVersion, _ = strconv.Atoi(headerParts[5])
	} else {
		// TODO return error
	}

	// Decode fields
	fieldsParts := strings.Split(fieldsStr, ",")
	for _, f := range fieldsParts {
		container.Fields = append(container.Fields, Field{Value: f})
	}

	return container
}

func encode() {
	fmt.Println("Encode")
}
