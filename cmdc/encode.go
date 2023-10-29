package cmdc

import (
	"strconv"

	"github.com/matrixxsoftware/go-mdd/mdd"
	log "github.com/sirupsen/logrus"
)

func Encode(containers *mdd.Containers) ([]byte, error) {

	// TODO handle multiple containers

	return encodeContainer(&containers.Containers[0])
}

func encodeContainer(container *mdd.Container) ([]byte, error) {
	log.Debugf("Encoding %+v\n", container)

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
		data = append(data, fields[0].Data...)
		// Remaining fields
		for i := 1; i < len(fields); i++ {
			data = append(data, ',')
			data = append(data, fields[i].Data...)
		}
	}
	data = append(data, ']')

	return data, nil
}
