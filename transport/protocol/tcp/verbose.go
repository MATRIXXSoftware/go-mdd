package tcp

import (
	"fmt"
	"net"
	"strings"
)

func connStr(conn net.Conn) string {
	remote := conn.RemoteAddr().String()
	local := conn.LocalAddr().String()
	return fmt.Sprintf("[%s->%s]", local, remote)
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
