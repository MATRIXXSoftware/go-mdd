package mdd

import (
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	ln net.Listener
}

// TODO make this configurable
const numWorkers = 10

func NewServer(addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{ln: ln}, nil
}

func (s *Server) Listen() error {
	jobs := make(chan net.Conn, 100)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go handleJobs(jobs, &wg)
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

func handleJobs(jobs chan net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for conn := range jobs {
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		request := &Containers{}
		err := request.Decode(conn)
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Received Request: %+v", request)

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

		err = response.Encode(conn)
		if err != nil {
			log.Panic(err)
		}
	}
}
