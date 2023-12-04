package tcp

import (
	"encoding/binary"
	"io"
)

func Encode(w io.Writer, encoded []byte) error {

	len := uint32(len(encoded))

	len += 4

	if err := binary.Write(w, binary.BigEndian, len); err != nil {
		return err
	}

	if _, err := w.Write(encoded); err != nil {
		return err
	}

	return nil
}

func Decode(r io.Reader) ([]byte, error) {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}

	len -= 4

	payload := make([]byte, len)

	_, err := io.ReadFull(r, payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
