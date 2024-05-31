package tcp

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

func Write(conn net.Conn, encoded []byte) error {
	len := uint32(len(encoded))

	if log.IsLevelEnabled(log.TraceLevel) {
		log.Tracef("%s Writing %d bytes to TCP",
			connStr(conn),
			len)
	}

	len += 4

	if err := binary.Write(conn, binary.BigEndian, len); err != nil {
		return err
	}

	n, err := conn.Write(encoded)
	if err != nil {
		return err
	}

	if log.IsLevelEnabled(log.TraceLevel) {
		hexdump := PrettyHexDump(encoded)
		log.Tracef("%s Written %d bytes to TCP:\n%s",
			connStr(conn),
			n,
			hexdump)
	}

	return nil
}

func Read(ctx context.Context, conn net.Conn) ([]byte, error) {
	var len uint32
	if err := binary.Read(conn, binary.BigEndian, &len); err != nil {
		return nil, err
	}

	len -= 4

	if log.IsLevelEnabled(log.TraceLevel) {
		log.Tracef("%s Reading %d bytes from TCP",
			connStr(conn),
			len)
	}

	payload := make([]byte, len)

	type Result struct {
		n   int
		err error
	}
	c := make(chan Result, 1)

	go func() {
		n, err := io.ReadFull(conn, payload)
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
					log.Tracef("%s Partial data read %d bytes from TCP:\n%s",
						connStr(conn),
						n,
						PrettyHexDump(payload[:n]))
				}
			}
			return nil, err
		}
		if log.IsLevelEnabled(log.TraceLevel) {
			hexdump := PrettyHexDump(payload)
			log.Tracef("%s Read %d bytes from TCP:\n%s",
				connStr(conn),
				n,
				hexdump)
		}

		return payload, nil
	}
}

func connStr(conn net.Conn) string {
	remote := conn.RemoteAddr().String()
	local := conn.LocalAddr().String()
	return fmt.Sprintf("[%s/%s]", local, remote)
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
