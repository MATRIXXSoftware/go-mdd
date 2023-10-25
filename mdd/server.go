package mdd

import (
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	ln    net.Listener
	codec Codec
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

		log.Printf("Received Request: %+v", requestTransport.Containers)

		// TODO add callback

		// Dummy Response for now
		response := Containers{
			Containers: []Container{
				{
					Header: Header{
						Version:       1,
						TotalField:    2,
						Depth:         0,
						Key:           88,
						SchemaVersion: 5222,
						ExtVersion:    2,
					},
					Fields: []Field{{Value: "Ok"}, {Value: "0"}},
				},
			},
		}

		responseTransport := Transport{
			Containers: &response,
			Codec:      s.codec,
		}

		err = responseTransport.Encode(conn)
		if err != nil {
			log.Panic(err)
		}
	}
}
