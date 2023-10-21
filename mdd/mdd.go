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

func Decode2(data []byte) (Container, error) {
	fmt.Println("Decoding ", data)
	var container Container
	var headerData []byte

	// First char must be '<'
	if data[0] != '<' {
		return container, errors.New("Invalid cMDC header, first char must be '<'")
	}
	headerStarted := true

	// Start from second char
	var idx int
	for idx = 1; idx < len(data); idx++ {
		c := data[idx]
		//log.Printf("i:%-3d c:%c\n", idx, c)
		if headerStarted && c == '>' {
			headerStarted = false
			headerData = data[1:idx]
			break
		}
	}

	// log.Printf("headerString: %s\n", headerData)
	header, err := decodeHeader2(headerData)
	if err != nil {
		return container, err
	}

	container.Header = header

	// Decode Body
	idx++
	if idx >= len(data) {
		return container, errors.New("Invalid cMDC body, no body")
	}

	// log.Printf("idx: %d %c\n", idx, data[idx])

	// First char must be '['
	if data[idx] != '[' {
		return container, errors.New("Invalid cMDC body, first char must be '['")
	}

	return container, nil
}

func decodeHeader2(data []byte) (Header, error) {
	var header Header
	// TODO do not use string conversion
	headerParts := strings.Split(string(data), ",")
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

func Encode(container Container) (string, error) {
	fmt.Println("Encoding ", container)
	var data string

	// Encode header
	header := container.Header
	data += fmt.Sprintf("<%d,%d,%d,%d,%d,%d>", header.Version, header.TotalField, header.Depth, header.Key, header.SchemaVersion, header.ExtVersion)

	// Encode fields
	for _, f := range container.Fields {
		data += fmt.Sprintf("%s,", f.Value)
	}

	// Remove last comma
	data = data[:len(data)-1]

	return data, nil
}
