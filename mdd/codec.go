package mdd

import (
	"encoding/binary"
	"io"
)

type Codec interface {
	Decode([]byte) (*Containers, error)
	Encode(*Containers) ([]byte, error)
}

func Encode(w io.Writer, codec Codec, containers *Containers) error {
	encoded, err := codec.Encode(containers)
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

func Decode(r io.Reader, codec Codec) (*Containers, error) {
	var len uint32
	if err := binary.Read(r, binary.LittleEndian, &len); err != nil {
		return nil, err
	}

	payload := make([]byte, len)

	_, err := io.ReadFull(r, payload)
	if err != nil {
		return nil, err
	}

	decoded, err := codec.Decode(payload)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}
