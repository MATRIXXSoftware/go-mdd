package cmdc

import (
	"strconv"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func (cmdc *Cmdc) encodeContainers(containers *mdd.Containers) ([]byte, error) {

	var data []byte
	for i := range containers.Containers {
		container := &containers.Containers[i]
		containerData, err := cmdc.encodeContainer(container)
		if err != nil {
			return nil, err
		}
		data = append(data, containerData...)
	}

	return data, nil
}

func (cmdc *Cmdc) encodeContainer(container *mdd.Container) ([]byte, error) {

	var data []byte

	headerData, err := cmdc.encodeHeader(&container.Header)
	if err != nil {
		return nil, err
	}
	data = append(data, headerData...)

	bodyData, err := cmdc.encodeBody(container.Fields)
	if err != nil {
		return nil, err
	}
	data = append(data, bodyData...)

	return data, nil
}

func (cmdc *Cmdc) encodeHeader(header *mdd.Header) ([]byte, error) {
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

func (cmdc *Cmdc) encodeBody(fields []mdd.Field) ([]byte, error) {
	// Pre-allocate a slice of bytes for better performance
	estimatedLen := len(fields) + 2
	for i := range fields {
		f := &fields[i]
		estimatedLen += len(f.Data)
	}
	data := make([]byte, 0, estimatedLen)

	data = append(data, '[')
	if len(fields) != 0 {
		// First field
		fieldData, err := cmdc.EncodeField(&fields[0])
		if err != nil {
			return nil, err
		}
		data = append(data, fieldData...)

		// Remaining fields
		for i := 1; i < len(fields); i++ {
			data = append(data, ',')
			fieldData, err := cmdc.EncodeField(&fields[i])
			if err != nil {
				return nil, err
			}
			data = append(data, fieldData...)
		}
	}
	data = append(data, ']')

	return data, nil
}
