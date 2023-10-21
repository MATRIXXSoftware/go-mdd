package mdd

import (
	"errors"
	"fmt"
	"log"
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

func Decode(data []byte) (Container, error) {
	log.Printf("Decoding %s\n", string(data))
	var container Container

	// Decode Header
	var headerData []byte

	// First char must be '<'
	if data[0] != '<' {
		return container, errors.New("Invalid cMDC header, first char must be '<'")
	}

	// Start from second char
	var idx int
	for idx = 1; idx < len(data); idx++ {
		c := data[idx]
		// log.Printf("i:%-3d c:%c\n", idx, c)
		if c == '>' {
			headerData = data[1:idx]
			break
		}
	}
	// log.Printf("headerString: %s\n", headerData)
	header, err := decodeHeader(headerData)
	if err != nil {
		return container, err
	}
	container.Header = header

	// Decode Body
	var bodyData []byte

	idx++
	if idx >= len(data) {
		return container, errors.New("Invalid cMDC body, no body")
	}

	// log.Printf("idx: %d %c\n", idx, data[idx])

	// First char must be '['
	if data[idx] != '[' {
		return container, errors.New("Invalid cMDC body, first char must be '['")
	}
	bodyStartIdx := idx
	for ; idx < len(data); idx++ {
		c := data[idx]
		// log.Printf("i:%-3d c:%c\n", idx, c)
		if c == ']' {
			bodyData = data[bodyStartIdx+1 : idx]
			break
		}
	}

	// log.Printf("bodyString: %s\n", bodyData)

	fields, err := decodeBody(bodyData)
	if err != nil {
		return container, err
	}

	container.Fields = fields

	return container, nil
}

func decodeBody(data []byte) ([]Field, error) {
	log.Printf("Decoding body '%s'\n", string(data))

	var fields []Field

	mark := 0
	i := 0

	for ; i < len(data); i++ {
		c := data[i]
		// log.Printf("c: %c\n", c)
		if c == ',' {
			fieldData := data[mark:i]
			// log.Printf("fieldData: %s\n", fieldData)
			mark = i + 1
			field := Field{Value: string(fieldData)}
			fields = append(fields, field)
		}
	}
	// last field
	fieldData := data[mark:i]
	// log.Printf("fieldData: %s\n", fieldData)
	field := Field{Value: string(fieldData)}
	fields = append(fields, field)

	return fields, nil
}

func decodeHeader(data []byte) (Header, error) {
	log.Printf("Decoding header '%s'\n", string(data))

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
