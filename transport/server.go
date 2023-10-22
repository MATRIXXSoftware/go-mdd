package transport

import (
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

type MddServer struct {
	ln net.Listener
}

// TODO make this configurable
const numWorkers = 10

func NewServer(addr string) (*MddServer, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &MddServer{ln: ln}, nil
}

func (s *MddServer) Listen() error {
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

	msg := &ExampleMessage{}
	err := msg.Decode(conn)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Received message: %+v", msg)

	err = msg.Encode(conn)
	if err != nil {
		log.Panic(err)
	}
}
