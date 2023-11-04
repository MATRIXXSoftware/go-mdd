package tcp

import (
	"io"
	"net"

	"github.com/matrixxsoftware/go-mdd/mdd"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	ln      net.Listener
	codec   mdd.Codec
	handler func(*mdd.Containers) *mdd.Containers
}

func NewServer(addr string, codec mdd.Codec) (*Server, error) {

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		ln:    ln,
		codec: codec,
	}, nil
}

func (s *Server) Listen() error {
	for {
		// Accept a connection
		conn, err := s.ln.Accept()
		if err != nil {
			opErr, ok := err.(*net.OpError)
			if ok && opErr.Op == "accept" {
				log.Infof("Server shutting down")
				return nil
			}
			return err
		}
		// Spawn a new Goroutine for each incoming connection
		go s.handleConnection(conn)
	}
}

func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) Handler(handler func(*mdd.Containers) *mdd.Containers) {
	s.handler = handler
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Handle message synchronously
	for {
		request, err := Decode(conn, s.codec)
		if err != nil {
			if err == io.EOF {
				log.Infof("Connection closed")
				return
			}
			log.Panic(err)
		}

		response := s.handler(request)

		err = Encode(conn, s.codec, response)
		if err != nil {
			log.Panic(err)
		}
	}
}
