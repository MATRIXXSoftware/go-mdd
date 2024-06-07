package tcp

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
	log "github.com/sirupsen/logrus"
)

type ServerTransport struct {
	ln      net.Listener
	handler func(*mdd.Containers) (*mdd.Containers, error)
	Codec   mdd.Codec
}

func NewServerTransport(addr string, codec mdd.Codec) (*ServerTransport, error) {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &ServerTransport{
		ln:    ln,
		Codec: codec,
	}, nil
}

func (s *ServerTransport) Listen() error {
	for {
		// Accept a connection
		conn, err := s.ln.Accept()
		if err != nil {
			opErr, ok := err.(*net.OpError)
			if ok && opErr.Op == "accept" {
				log.Infof("ServerTransport shutting down")
				return nil
			}
			return err
		}
		// Spawn a new Goroutine for each incoming connection
		go s.handleConnection(conn)
	}
}

func (s *ServerTransport) Close() error {
	return s.ln.Close()
}

func (s *ServerTransport) Handler(handler func(*mdd.Containers) (*mdd.Containers, error)) {
	s.handler = handler
}

func (s *ServerTransport) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Handle message synchronously
	for {
		reqBody, err := Read(conn)
		if err != nil {
			if err == io.EOF {
				log.Infof("%s Connection closed", connStr(conn))
			} else if err == io.ErrUnexpectedEOF {
				log.Errorf("%s Connection closed unexpectedly", connStr(conn))
			} else if err == context.DeadlineExceeded {
				log.Errorf("%s Timeout reading from connection", connStr(conn))
			} else {
				log.Errorf("%s Error reading from connection: %s", connStr(conn), err)
			}
			return
		}

		respBody, err := s.processMessage(reqBody)
		if err != nil {
			log.Errorf("%s %s", connStr(conn), err)
			return
		}

		err = Write(conn, respBody)
		if err != nil {
			log.Errorf("%s %s", connStr(conn), err)
			return
		}
	}
}

func (s *ServerTransport) processMessage(reqBody []byte) ([]byte, error) {
	req, err := s.Codec.Decode(reqBody)
	if err != nil {
		return nil, err
	}

	// hopId, err := extractHopId(req)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("HopId: %d\n", hopId)

	response, err := s.handler(req)
	if err != nil {
		return nil, err
	}

	respBody, err := s.Codec.Encode(response)
	if err != nil {
		return nil, err
	}
	return respBody, nil
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
