package cmdc

import (
	"strconv"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
)

func Encode(containers *mdd.Containers) ([]byte, error) {

	var data []byte
	for _, container := range containers.Containers {
		containerData, err := encodeContainer(&container)
		if err != nil {
			return nil, err
		}
		data = append(data, containerData...)
	}

	return data, nil
}

func encodeContainer(container *mdd.Container) ([]byte, error) {

	var data []byte

	headerData, err := encodeHeader(&container.Header)
	if err != nil {
		return nil, err
	}
	data = append(data, headerData...)

	bodyData, err := encodeBody(container.Fields)
	if err != nil {
		return nil, err
	}
	data = append(data, bodyData...)

	return data, nil
}

func encodeHeader(header *mdd.Header) ([]byte, error) {
	// Estimate size:
	// version          4
	// totalField       1
	// depth            1
	// key              7
	// schemaVersion    4
	// extVersion       3
	// commas           6
	// angle brackets   2
	estimatedSize := 4 + 1 + 1 + 7 + 4 + 3 + 6 + 2
	b := make([]byte, 0, estimatedSize)

	b = append(b, '<')
	b = strconv.AppendInt(b, int64(header.Version), 10)
	b = append(b, ',')
	b = strconv.AppendInt(b, int64(header.TotalField), 10)
	b = append(b, ',')
	b = strconv.AppendInt(b, int64(header.Depth), 10)
	b = append(b, ',')
	b = strconv.AppendInt(b, int64(header.Key), 10)
	b = append(b, ',')
	b = strconv.AppendInt(b, int64(header.SchemaVersion), 10)
	b = append(b, ',')
	b = strconv.AppendInt(b, int64(header.ExtVersion), 10)
	b = append(b, '>')

	return b, nil
}

func encodeBody(fields []mdd.Field) ([]byte, error) {
	// Pre-allocate a slice of bytes for better performance
	estimatedLen := len(fields) + 2
	for _, f := range fields {
		estimatedLen += len(f.Data)
	}
	data := make([]byte, 0, estimatedLen)

	data = append(data, '[')
	if len(fields) != 0 {
		// First field
		fieldData, err := encodeField(fields[0])
		if err != nil {
			return nil, err
		}
		data = append(data, fieldData...)

		// Remaining fields
		for i := 1; i < len(fields); i++ {
			data = append(data, ',')
			fieldData, err := encodeField(fields[i])
			if err != nil {
				return nil, err
			}
			data = append(data, fieldData...)
		}
	}
	data = append(data, ']')

	return data, nil
}

func encodeField(f mdd.Field) ([]byte, error) {
	// If the f has data, use it
	if len(f.Data) > 0 || f.Type == field.Unknown {
		return f.Data, nil
	}

	// Otherwise, encode the value
	switch f.Type {
	case field.Int32:
		v := f.Value.Int32()
		return []byte(strconv.Itoa(int(v))), nil

	case field.String:
		v := f.Value.String()
		data := make([]byte, 0, len(v)+6)
		data = append(data, '(')
		data = append(data, []byte(strconv.Itoa(len(v)))...)
		data = append(data, ':')
		data = append(data, []byte(v)...)
		data = append(data, ')')
		return data, nil

	case field.Struct:
		containers := f.Value.Struct()
		return encodeContainer(&containers.Containers[0])

	// TODO support other types

	default:
		return f.Data, nil
	}
}
