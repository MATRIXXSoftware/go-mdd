package transport

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
)

type ExampleMessage struct {
	Field1 uint32
	Field2 uint32
}

func (m *ExampleMessage) Decode(r io.Reader) error {
	return binary.Read(r, binary.LittleEndian, m)
}

func (m *ExampleMessage) Encode(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, m)
}

const numWorkers = 10

func MddServer(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	jobs := make(chan net.Conn, 100)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go handleJobs(jobs, &wg)
	}

	for {
		conn, err := ln.Accept()
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
