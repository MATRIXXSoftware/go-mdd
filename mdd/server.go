package mdd

import (
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	ln      net.Listener
	codec   Codec
	handler func(*Containers) *Containers
}

// TODO make this configurable
const numWorkers = 10

func NewServer(addr string, codec Codec) (*Server, error) {

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
	jobs := make(chan net.Conn, 100)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.handleJobs(jobs, &wg)
	}

	for {
		conn, err := s.ln.Accept()
		if err != nil {
			break
		}
		jobs <- conn
	}

	wg.Wait()
	return nil
}

func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) Handler(handler func(*Containers) *Containers) {
	s.handler = handler
}

func (s *Server) handleJobs(jobs chan net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for conn := range jobs {
		s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		requestTransport := Transport{
			Containers: &Containers{},
			Codec:      s.codec,
		}
		err := requestTransport.Decode(conn)
		if err != nil {
			log.Panic(err)
		}

		response := s.handler(requestTransport.Containers)

		responseTransport := Transport{
			Containers: response,
			Codec:      s.codec,
		}

		err = responseTransport.Encode(conn)
		if err != nil {
			log.Panic(err)
		}
	}
}
