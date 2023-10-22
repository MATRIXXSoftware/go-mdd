package transport

import (
	"encoding/binary"
	"io"
)

// TODO replace this with MDD Container
type ExampleMessage struct {
	Field1 uint32
	Field2 uint32
}

func (m *ExampleMessage) Decode(r io.Reader) error {
	return binary.Read(r, binary.LittleEndian, m)
}

func (m *ExampleMessage) Encode(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, m)
}
