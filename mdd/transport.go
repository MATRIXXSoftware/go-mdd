package mdd

import (
	"encoding/binary"
	"io"
)

type Transport struct {
	Codec      Codec
	Containers *Containers
}

func (t *Transport) Encode(w io.Writer) error {
	encodedStr, err := t.Codec.Encode(t.Containers)
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

func (t *Transport) Decode(r io.Reader) error {
	var payloadLen uint32
	if err := binary.Read(r, binary.LittleEndian, &payloadLen); err != nil {
		return err
	}

	payload := make([]byte, payloadLen)

	_, err := io.ReadFull(r, payload)
	if err != nil {
		return err
	}

	decodedContainers, err := t.Codec.Decode(payload)
	if err != nil {
		return err
	}

	t.Containers = decodedContainers

	return nil
}
