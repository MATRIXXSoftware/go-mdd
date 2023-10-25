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
	encoded, err := t.Codec.Encode(t.Containers)
	if err != nil {
		return err
	}

	len := uint32(len(encoded))

	if err := binary.Write(w, binary.LittleEndian, len); err != nil {
		return err
	}

	_, err = w.Write(encoded)
	if err != nil {
		return err
	}

	return nil
}

func (t *Transport) Decode(r io.Reader) error {
	var len uint32
	if err := binary.Read(r, binary.LittleEndian, &len); err != nil {
		return err
	}

	payload := make([]byte, len)

	_, err := io.ReadFull(r, payload)
	if err != nil {
		return err
	}

	decoded, err := t.Codec.Decode(payload)
	if err != nil {
		return err
	}

	t.Containers = decoded

	return nil
}
