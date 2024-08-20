package tcp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
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

func Read(conn net.Conn) ([]byte, error) {
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))

	var len uint32
	if err := binary.Read(conn, binary.BigEndian, &len); err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return nil, nil // timeout is expected due to deadline
		}
		return nil, err
	}

	len -= 4

	if log.IsLevelEnabled(log.TraceLevel) {
		log.Tracef("%s Reading %d bytes from TCP",
			connStr(conn),
			len)
	}

	payload := make([]byte, len)

	n, err := io.ReadFull(conn, payload)
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

func extractHopId(containers *mdd.Containers) (uint32, error) {
	// Get MtxMsg container (key 93)
	mtxMsg := containers.GetContainer(93)
	if mtxMsg == nil {
		return 0, fmt.Errorf("container MtxMsg is missing")
	}

	// Assume no changes to the position of hopId field
	f := mtxMsg.GetField(14)

	if f.Data == nil {
		return 0, fmt.Errorf("hopId field is missing")
	}

	// Copy the field data to a new field
	hopIdField := mdd.Field{
		Data:  f.Data,
		Type:  field.UInt32,
		Codec: f.Codec,
	}

	// Get the value of the field
	hopId, err := hopIdField.GetValue()
	if err != nil {
		return 0, err
	}

	return hopId.(uint32), nil
}

func injectHopId(containers *mdd.Containers, hopId uint32) error {
	// Get MtxMsg container (key 93)
	mtxMsg := containers.GetContainer(93)
	if mtxMsg == nil {
		return fmt.Errorf("container MtxMsg is missing")
	}

	// Assume no changes to the position of hopId field
	mtxMsg.SetField(14, &mdd.Field{
		Value: hopId,
		Type:  field.UInt32,
	})

	return nil
}
