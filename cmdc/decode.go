package cmdc

import (
	"errors"
	"fmt"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func (cmdc *Cmdc) decodeContainers(data []byte) (*mdd.Containers, error) {

	var containers mdd.Containers
	var idx int
	slice := data

	for idx < len(slice) {
		container, idx, err := cmdc.decodeContainer(slice)
		if err != nil {
			return nil, err
		}

		containers.Containers = append(containers.Containers, *container)
		slice = slice[idx:]
	}

	return &containers, nil
}

func (cmdc *Cmdc) decodeContainer(data []byte) (*mdd.Container, int, error) {
	var container mdd.Container
	idx := 0

	// Decode Header
	header, offset, err := cmdc.decodeHeader(data[idx:])
	if err != nil {
		return nil, idx, err
	}
	container.Header = header
	idx += offset

	// Decode Body
	fields, offset, err := cmdc.decodeBody(data[idx:])
	if err != nil {
		return &container, idx, err
	}
	container.Fields = fields
	idx += offset

	// Load Definition
	if cmdc.dict != nil {
		definition, ok := cmdc.dict.Get(container.Header.Key)
		if ok {
			container.LoadDefinition(definition)
		}
	}

	return &container, idx, nil
}

func (cmdc *Cmdc) decodeBody(data []byte) ([]mdd.Field, int, error) {

	var fields []mdd.Field

	// Check if there is a body
	if len(data) <= 0 {
		return nil, 0, errors.New("invalid cMDC body, no body")
	}

	// First char following a header must be '['
	if data[0] != '[' {
		return nil, 0, errors.New("invalid cMDC body, first character must be '['")
	}

	idx := 1
	mark := idx
	roundMark := 0

	square := 1
	angle := 0
	round := 0
	curly := 0

	isMulti := false
	isContainer := false

	complete := false

	for ; idx < len(data); idx++ {
		c := data[idx]

		if round != 0 {
			if c == ')' {
				round--
			} else if roundMark == 0 {
				return nil, idx, errors.New("invalid cMDC body, mismatch string length")
			} else if c == ':' {
				temp := data[roundMark+1 : idx]
				len, err := bytesToInt(temp)
				if err != nil {
					return nil, idx, errors.New("invalid cMDC body, invalid string length")
				}
				// reset round mark
				roundMark = 0
				// skip the string field
				idx += len
			} else if c < '0' || c > '9' {
				return nil, idx, errors.New("invalid character '" + string(c) + "', numeric expected for string length")
			}
			continue
		}

		switch c {
		case '(':
			roundMark = idx
			round++
		case '[':
			square++
		case ']':
			square--
		case '<':
			isContainer = true
			angle++
		case '>':
			angle--
		case '{':
			if square == 1 {
				isMulti = true
			}
			curly++
		case '}':
			curly--
		case ',':
			if square == 1 && angle == 0 && curly == 0 {
				// Extract field
				fieldData := data[mark:idx]
				mark = idx + 1
				field := mdd.Field{
					Data:        fieldData,
					IsMulti:     isMulti,
					IsContainer: isContainer,
					IsNull:      len(fieldData) == 0,
					Codec:       cmdc,
					Value:       nil,
				}

				fields = append(fields, field)
				isMulti = false
				isContainer = false
			}
		}

		// End of body
		if square == 0 {
			complete = true
			idx++
			break
		}
	}

	if complete == false {
		return nil, idx, errors.New("invalid cMDC body, no end of body")
	}

	// Extract last field
	fieldData := data[mark : idx-1]
	field := mdd.Field{
		Data:        fieldData,
		IsMulti:     isMulti,
		IsContainer: isContainer,
		IsNull:      len(fieldData) == 0,
		Codec:       cmdc,
		Value:       nil,
	}
	fields = append(fields, field)

	return fields, idx, nil
}

func (cmdc *Cmdc) decodeHeader(data []byte) (mdd.Header, int, error) {

	var header mdd.Header

	// Check if there is a header
	if len(data) <= 0 {
		return header, 0, errors.New("invalid cMDC header, no header")
	}

	// First char must be '<'
	if data[0] != '<' {
		return header, 0, errors.New("invalid cMDC header, first character must be '<'")
	}

	idx := 1
	mark := idx
	fieldNumber := 0
	complete := false
	for ; idx < len(data); idx++ {
		c := data[idx]
		if c == '>' {
			complete = true
			idx++
			break
		} else if c == ',' {
			fieldData := data[mark:idx]
			mark = idx + 1
			v, err := bytesToInt(fieldData)
			if err != nil {
				return header, idx, errors.New("invalid cMDC header field '" + string(fieldData) + "', numeric expected")
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
			return header, idx, errors.New("invalid cMDC character '" + string(c) + "' in header, numeric expected")
		}
	}

	if complete == false {
		return header, idx, errors.New("invalid cMDC header, no end of header")
	}

	// last field
	fieldData := data[mark : idx-1]
	v, err := bytesToInt(fieldData)
	if err != nil {
		return header, idx, err
	}

	if fieldNumber != 5 {
		return header, idx, errors.New("invalid cMDC header, 6 fields expected")
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
			return 0, fmt.Errorf("invalid character in byte slice: %c", b[i])
		}
		result += int(b[i]-'0') * multiplier
		multiplier *= 10
	}

	if isNegative {
		result = -result
	}

	return result, nil
}
