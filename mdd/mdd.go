package mdd

import (
	"errors"
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

func Decode(data string) (Container, error) {
	fmt.Println("Decoding ", data)
	var container Container

	// Extract header and fields substrings
	headerEnd := strings.Index(data, ">")
	fieldsStart := strings.Index(data, "[")

	headerStr := data[1:headerEnd]
	fieldsStr := data[fieldsStart+1 : len(data)-1]

	// Decode Header
	header, err := decodeHeader(headerStr)
	if err != nil {
		return container, err
	}

	container.Header = header

	// Decode fields
	fieldsParts := strings.Split(fieldsStr, ",")
	for _, f := range fieldsParts {
		container.Fields = append(container.Fields, Field{Value: f})
	}

	return container, nil
}

func decodeHeader(data string) (Header, error) {
	var header Header
	headerParts := strings.Split(data, ",")
	if len(headerParts) != 6 {
		return header, errors.New("Invalid cMDC header")
	}
	header.Version, _ = strconv.Atoi(headerParts[0])
	header.TotalField, _ = strconv.Atoi(headerParts[1])
	header.Depth, _ = strconv.Atoi(headerParts[2])
	header.Key, _ = strconv.Atoi(headerParts[3])
	header.SchemaVersion, _ = strconv.Atoi(headerParts[4])
	header.ExtVersion, _ = strconv.Atoi(headerParts[5])
	return header, nil
}

func Encode() {
	fmt.Println("Encode")
}
