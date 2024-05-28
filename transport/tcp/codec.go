package tcp

import (
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Encode(w io.Writer, encoded []byte) error {

	len := uint32(len(encoded))

	log.Tracef("Writing to TCP stream with MDC payload of length: %d", len)

	len += 4

	if err := binary.Write(w, binary.BigEndian, len); err != nil {
		return err
	}

	if _, err := w.Write(encoded); err != nil {
		return err
	}

	if log.IsLevelEnabled(log.TraceLevel) {
		hexdump := PrettyHexDump(encoded)
		log.Trace("MDC payload written to TCP stream:")
		fmt.Printf(hexdump)
	}

	return nil
}

func Decode(r io.Reader) ([]byte, error) {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}

	len -= 4

	log.Tracef("Expecting to read from TCP stream MDC payload of length: %d", len)

	payload := make([]byte, len)

	_, err := io.ReadFull(r, payload)
	if err != nil {
		return nil, err
	}

	if log.IsLevelEnabled(log.TraceLevel) {
		hexdump := PrettyHexDump(payload)
		log.Trace("MDC payload read from TCP stream:")
		fmt.Printf(hexdump)
	}

	return payload, nil
}

func PrettyHexDump(data []byte) string {
	const bytesPerLine = 16
	var builder strings.Builder

	for i := 0; i < len(data); i += bytesPerLine {
		// Print the offset
		builder.WriteString(fmt.Sprintf("%08x  ", i))

		// Print the hex codes
		for j := 0; j < bytesPerLine; j++ {
			if i+j < len(data) {
				builder.WriteString(fmt.Sprintf("%02x ", data[i+j]))
			} else {
				builder.WriteString("   ")
			}
		}

		// Print the separator
		builder.WriteString(" ")

		// Print the ASCII representation
		for j := 0; j < bytesPerLine; j++ {
			if i+j < len(data) {
				b := data[i+j]
				if b >= 32 && b <= 126 {
					builder.WriteByte(b)
				} else {
					builder.WriteByte('.')
				}
			}
		}

		builder.WriteString("\n")
	}

	return builder.String()
}
