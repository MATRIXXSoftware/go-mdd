package tcp

import (
	"context"
	"io"
	"net"
	"time"

	// "time"

	log "github.com/sirupsen/logrus"
)

type ServerTransport struct {
	ln      net.Listener
	handler func([]byte) ([]byte, error)
}

func NewServerTransport(addr string) (*ServerTransport, error) {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &ServerTransport{
		ln: ln,
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

func (s *ServerTransport) Handler(handler func([]byte) ([]byte, error)) {
	s.handler = handler
}

func (s *ServerTransport) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Handle message synchronously
	for {
		// Hard code 3 second for now. Make it configurable later
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		request, err := Read(ctx, conn)
		if err != nil {
			if err == io.EOF {
				log.Infof("Connection closed")
			} else if err == io.ErrUnexpectedEOF {
				log.Errorf("Connection closed unexpectedly")
			} else if err == context.DeadlineExceeded {
				log.Errorf("Timeout reading from connection")
			} else {
				log.Errorf("Error reading from connection: %s", err)
			}
			return
		}

		response, err := s.handler(request)
		if err != nil {
			log.Errorf("%s", err)
			return
		}

		err = Write(conn, response)
		if err != nil {
			log.Errorf("%s", err)
			return
		}
	}
}
