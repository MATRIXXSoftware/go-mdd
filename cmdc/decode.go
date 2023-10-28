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
	idx := 0

	// Decode Header
	header, offset, err := decodeHeader(data[idx:])
	if err != nil {
		return nil, idx, err
	}
	container.Header = header
	idx += offset

	// Decode Body
	fields, offset, err := decodeBody(data[idx:])
	if err != nil {
		return &container, idx, err
	}
	container.Fields = fields
	idx += offset

	return &container, idx, nil
}

func decodeBody(data []byte) ([]mdd.Field, int, error) {

	var fields []mdd.Field

	// Check if there is a body
	if len(data) <= 0 {
		return nil, 0, errors.New("Invalid cMDC body, no body")
	}

	// First char following a header must be '['
	if data[0] != '[' {
		return nil, 0, errors.New("Invalid cMDC body, first char must be '['")
	}

	idx := 1
	mark := idx
	roundMark := 0

	square := 1
	angle := 0
	round := 0

	for ; idx < len(data); idx++ {
		c := data[idx]

		if round != 0 {
			if c == ':' {
				temp := data[roundMark+1 : idx]
				len, err := bytesToInt(temp)
				if err != nil {
					panic("Invalid string length")
				}
				// skip the string field
				idx += len
			} else if c == ')' {
				round--
			} else if c < '0' || c > '9' {
				panic("Invalid character, numeric expected for string length")
			}
			continue
		}

		switch c {
		case '(':
			round++
			roundMark = idx
		case '[':
			square++
		case ']':
			square--
		case '<':
			angle++
		case '>':
			angle--
		case ',':
			if square == 1 && angle == 0 {
				// Extract field
				fieldData := data[mark:idx]
				mark = idx + 1
				field := mdd.Field{
					Data: fieldData,
				}
				fields = append(fields, field)
			}
		}

		// End of body
		if square == 0 {
			idx++
			break
		}
	}

	// Extract last field
	fieldData := data[mark : idx-1]
	field := mdd.Field{
		Data: fieldData,
	}
	fields = append(fields, field)

	return fields, idx, nil
}

func decodeHeader(data []byte) (mdd.Header, int, error) {

	var header mdd.Header

	// Check if there is a header
	if len(data) <= 0 {
		return header, 0, errors.New("Invalid cMDC header, no header")
	}

	// First char must be '<'
	if data[0] != '<' {
		return header, 0, errors.New("Invalid cMDC header, first char must be '<'")
	}

	idx := 1
	mark := idx
	fieldNumber := 0
	for ; idx < len(data); idx++ {
		c := data[idx]
		if c == '>' {
			idx++
			break
		} else if c == ',' {
			fieldData := data[mark:idx]
			mark = idx + 1
			v, err := bytesToInt(fieldData)
			if err != nil {
				return header, idx, err
			}
			switch fieldNumber {
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
			fieldNumber++
		} else if (c < '0' || c > '9') && c != '-' {
			return header, idx, errors.New("Invalid cMDC header, '" + string(c) + "' numeric expected")
		}
	}

	// last field
	fieldData := data[mark : idx-1]
	v, err := bytesToInt(fieldData)
	if err != nil {
		return header, idx, err
	}

	if fieldNumber != 5 {
		return header, idx, errors.New("Invalid cMDC header, 6 fields expected")
	}

	header.ExtVersion = v
	return header, idx, nil
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
