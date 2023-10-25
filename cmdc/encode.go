package cmdc

import (
	"fmt"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

func Encode(containers *mdd.Containers) ([]byte, error) {

	// TODO handle multiple containers

	return encodeContainer(&containers.Containers[0])
}

func encodeContainer(container *mdd.Container) ([]byte, error) {
	fmt.Println("Encoding ", container)
	var data string

	// Encode header
	header := container.Header
	data += fmt.Sprintf("<%d,%d,%d,%d,%d,%d>", header.Version, header.TotalField, header.Depth, header.Key, header.SchemaVersion, header.ExtVersion)

	// Encode fields
	data += "["
	for _, f := range container.Fields {
		data += fmt.Sprintf("%s,", f.Value)
	}

	// Remove last comma
	data = data[:len(data)-1]
	data += "]"

	return []byte(data), nil
}
