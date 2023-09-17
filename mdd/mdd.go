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

func Decode(data string) Container {
	fmt.Println("Decoding ", data)
	var container Container

	// Extract header and fields substrings
	headerEnd := strings.Index(data, ">")
	fieldsStart := strings.Index(data, "[")

	headerStr := data[1:headerEnd]
	fieldsStr := data[fieldsStart+1 : len(data)-1]

	// Decode Header
	container.Header = decodeHeader(headerStr)

	// Decode fields
	fieldsParts := strings.Split(fieldsStr, ",")
	for _, f := range fieldsParts {
		container.Fields = append(container.Fields, Field{Value: f})
	}

	return container
}

func decodeHeader(data string) Header {
	var header Header
	headerParts := strings.Split(data, ",")
	if len(headerParts) == 6 {
		header.Version, _ = strconv.Atoi(headerParts[0])
		header.TotalField, _ = strconv.Atoi(headerParts[1])
		header.Depth, _ = strconv.Atoi(headerParts[2])
		header.Key, _ = strconv.Atoi(headerParts[3])
		header.SchemaVersion, _ = strconv.Atoi(headerParts[4])
		header.ExtVersion, _ = strconv.Atoi(headerParts[5])
	} else {
		// TODO return error
	}
	return header
}

func Encode() {
	fmt.Println("Encode")
}
