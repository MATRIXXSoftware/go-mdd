package core

import (
	"errors"
	"fmt"
	"log"
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
		if c == '>' {
			headerData = data[1:idx]
			break
		}
	}
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

	// First char following a header must be '['
	if data[idx] != '[' {
		return container, errors.New("Invalid cMDC body, first char must be '['")
	}
	mark := idx
	for ; idx < len(data); idx++ {
		c := data[idx]
		if c == ']' {
			bodyData = data[mark+1 : idx]
			break
		}
	}

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
		if c == ',' {
			fieldData := data[mark:i]
			mark = i + 1
			field := Field{Value: string(fieldData)}
			fields = append(fields, field)
		}
	}
	// last field
	fieldData := data[mark:i]
	field := Field{Value: string(fieldData)}
	fields = append(fields, field)

	return fields, nil
}

func decodeHeader(data []byte) (Header, error) {
	log.Printf("Decoding header '%s'\n", string(data))
	var header Header
	mark := 0
	i := 0
	idx := 0
	for ; i < len(data); i++ {
		c := data[i]
		if c == ',' {
			fieldData := data[mark:i]
			mark = i + 1
			v, err := bytesToInt(fieldData)
			if err != nil {
				return header, err
			}
			switch idx {
			case 0:
				header.Version = v
			case 1:
				header.TotalField = v
			case 2:
				header.Depth = v
			case 3:
				header.Key = v
			case 4:
				header.SchemaVersion = v
			}
			idx++
		}
	}
	// last field
	fieldData := data[mark:i]
	v, err := bytesToInt(fieldData)
	if err != nil {
		return header, err
	}
	if idx != 5 {
		return header, errors.New("Invalid cMDC header, 6 fields expected")
	}
	header.ExtVersion = v
	return header, nil
}

func bytesToInt(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil // return zero if the slice is empty
	}

	result := 0
	multiplier := 1
	isNegative := false

	startIndex := 0
	if b[0] == '-' {
		isNegative = true
		startIndex = 1
	} else if b[0] == '+' {
		startIndex = 1
	}

	for i := len(b) - 1; i >= startIndex; i-- {
		if b[i] < '0' || b[i] > '9' {
			log.Fatalf("Invalid character in byte slice: %c", b[i])
			return 0, fmt.Errorf("Invalid character in byte slice: %c", b[i])
		}
		result += int(b[i]-'0') * multiplier
		multiplier *= 10
	}

	if isNegative {
		result = -result
	}

	return result, nil
}
