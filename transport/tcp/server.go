package tcp

import (
	"io"
	"net"

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
		request, err := Decode(conn)
		if err != nil {
			if err == io.EOF {
				log.Infof("Connection closed")
				return
			}
			log.Panic(err)
		}

		response, err := s.handler(request)
		if err != nil {
			log.Panic(err)
		}

		err = Encode(conn, response)
		if err != nil {
			log.Panic(err)
		}
	}
}
