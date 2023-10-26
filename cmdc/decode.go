package cmdc

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func Decode(data []byte) (*mdd.Containers, error) {
	log.Debugf("Decoding %s\n", string(data))

	var containers mdd.Containers
	var idx int
	slice := data

	for idx < len(slice) {
		container, idx, err := decodeContainer(slice)
		if err != nil {
			return nil, err
		}

		containers.Containers = append(containers.Containers, *container)
		slice = slice[idx:]
	}

	return &containers, nil
}

func decodeContainer(data []byte) (*mdd.Container, int, error) {
	var container mdd.Container

	// Decode Header
	var headerData []byte

	// First char must be '<'
	idx := 1
	if data[0] != '<' {
		return nil, idx, errors.New("Invalid cMDC header, first char must be '<'")
	}

	// Start from second char
	for ; idx < len(data); idx++ {
		c := data[idx]
		if c == '>' {
			headerData = data[1:idx]
			break
		}
	}
	header, err := decodeHeader(headerData)
	if err != nil {
		return nil, idx, err
	}
	container.Header = header

	// Decode Body
	var bodyData []byte

	idx++
	if idx >= len(data) {
		return nil, idx, errors.New("Invalid cMDC body, no body")
	}

	// First char following a header must be '['
	if data[idx] != '[' {
		return nil, idx, errors.New("Invalid cMDC body, first char must be '['")
	}

	mark := idx
	bracket := 0
	for ; idx < len(data); idx++ {
		c := data[idx]
		if c == '[' {
			bracket++
		} else if c == ']' {
			bracket--
			if bracket == 0 {
				bodyData = data[mark+1 : idx]
				idx++
				break
			}
		}
	}

	fields, err := decodeBody(bodyData)
	if err != nil {
		return &container, idx, err
	}

	container.Fields = fields

	return &container, idx, nil
}

func decodeBody(data []byte) ([]mdd.Field, error) {
	log.Debugf("Decoding body '%s'\n", string(data))

	var fields []mdd.Field

	mark := 0
	i := 0

	square := 0
	angle := 0
	for ; i < len(data); i++ {
		c := data[i]
		if c == '[' {
			square++
		} else if c == ']' {
			square--
		} else if c == '<' {
			angle++
		} else if c == '>' {
			angle--
		} else if square == 0 && angle == 0 && c == ',' {
			fieldData := data[mark:i]
			mark = i + 1
			field := mdd.Field{
				Data: fieldData,
			}
			fields = append(fields, field)
		}
	}
	// last field
	fieldData := data[mark:i]
	field := mdd.Field{
		Data: fieldData,
	}
	fields = append(fields, field)

	return fields, nil
}

func decodeHeader(data []byte) (mdd.Header, error) {
	log.Debugf("Decoding header '%s'\n", string(data))
	var header mdd.Header
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
