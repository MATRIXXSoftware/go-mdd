package cmdc

import (
	"strconv"

	"github.com/matrixxsoftware/go-mdd/mdd"
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
	str := "<" +
		strconv.Itoa(header.Version) + "," +
		strconv.Itoa(header.TotalField) + "," +
		strconv.Itoa(header.Depth) + "," +
		strconv.Itoa(header.Key) + "," +
		strconv.Itoa(header.SchemaVersion) + "," +
		strconv.Itoa(header.ExtVersion) +
		">"
	return []byte(str), nil
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

func encodeField(field mdd.Field) ([]byte, error) {
	if len(field.Data) > 0 || field.Type == mdd.Unknown {
		return field.Data, nil
	}

	switch field.Type {
	case mdd.Struct:
		containers := field.Value.(*mdd.Containers)
		return encodeContainer(&containers.Containers[0])
	}

	return field.Data, nil
}
