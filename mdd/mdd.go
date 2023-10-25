package mdd

import (
	"encoding/binary"
	"io"
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

type Containers struct {
	Containers []Container
}

func (c *Containers) GetContainer(key int) *Container {
	for _, container := range c.Containers {
		if container.Header.Key == key {
			return &container
		}
	}
	return nil
}

// TODO fix cyclic import
// TODO
func CmdcDecode(data []byte) (Containers, error) {
	response := Containers{
		Containers: []Container{
			{
				Header: Header{
					Version:       1,
					TotalField:    2,
					Depth:         0,
					Key:           88,
					SchemaVersion: 5222,
					ExtVersion:    2,
				},
				Fields: []Field{{Value: "Ok"}, {Value: "0"}},
			},
		},
	}
	return response, nil
}

// TODO
func CmdcEncode(container Containers) (string, error) {
	return "Dummy", nil
}


func (c *Containers) Encode(w io.Writer) error {
	encodedStr, err := CmdcEncode(*c)
	if err != nil {
		return err
	}

	encodedData := []byte(encodedStr)

	payloadLen := uint32(len(encodedData))

	if err := binary.Write(w, binary.LittleEndian, payloadLen); err != nil {
		return err
	}

	_, err = w.Write(encodedData)
	if err != nil {
		return err
	}

	return nil
}

func (c *Containers) Decode(r io.Reader) error {
	var payloadLen uint32
	if err := binary.Read(r, binary.LittleEndian, &payloadLen); err != nil {
		return err
	}

	payload := make([]byte, payloadLen)

	_, err := io.ReadFull(r, payload)
	if err != nil {
		return err
	}

	decodedContainers, err := CmdcDecode(payload)
	if err != nil {
		return err
	}

	*c = decodedContainers

	return nil
}
