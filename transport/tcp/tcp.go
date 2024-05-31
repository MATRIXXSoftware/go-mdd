package tcp

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Write(w io.Writer, encoded []byte) error {

	len := uint32(len(encoded))

	log.Tracef("Writing %d bytes to TCP", len)

	len += 4

	if err := binary.Write(w, binary.BigEndian, len); err != nil {
		return err
	}

	n, err := w.Write(encoded)
	if err != nil {
		return err
	}

	if log.IsLevelEnabled(log.TraceLevel) {
		hexdump := PrettyHexDump(encoded)
		log.Tracef("Written %d bytes to TCP:\n%s", n, hexdump)
	}

	return nil
}

func Read(ctx context.Context, r io.Reader) ([]byte, error) {
	var len uint32
	if err := binary.Read(r, binary.BigEndian, &len); err != nil {
		return nil, err
	}

	len -= 4

	log.Tracef("Reading %d bytes from TCP", len)

	payload := make([]byte, len)

	type Result struct {
		n   int
		err error
	}
	c := make(chan Result, 1)

	go func() {
		n, err := io.ReadFull(r, payload)
		c <- Result{n, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-c:
		n := res.n
		err := res.err
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				if log.IsLevelEnabled(log.TraceLevel) {
					log.Tracef("Partial data read %d bytes from TCP:\n%s",
						n,
						PrettyHexDump(payload[:n]))
				}
			}
			return nil, err
		}
		if log.IsLevelEnabled(log.TraceLevel) {
			hexdump := PrettyHexDump(payload)
			log.Tracef("Read %d bytes from TCP:\n%s", n, hexdump)
		}

		return payload, nil
	}
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
