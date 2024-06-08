package tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
	log "github.com/sirupsen/logrus"
)

type ClientTransport struct {
	conn     net.Conn
	Codec    mdd.Codec
	msgCache map[uint32]chan *mdd.Containers
	mu       sync.Mutex
}

func NewClientTransport(addr string, codec mdd.Codec) (*ClientTransport, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	c := &ClientTransport{
		conn:     conn,
		Codec:    codec,
		msgCache: make(map[uint32]chan *mdd.Containers),
	}

	// Processing goroutine
	go func() {
		for {
			// Read response from the server
			respBody, err := Read(c.conn)
			if err != nil {
				if err == io.EOF {
					log.Errorf("Connection closed by server")
					return
				} else if err == io.ErrUnexpectedEOF {
					log.Errorf("Connection closed unexpectedly by server")
					return
				}
				// TODO: handle error properly
				// panic(err)
				log.Errorf("Error reading from connection: %v", err)
				return
			}

			// Handle the response
			err = c.HandleResponse(respBody)
			if err != nil {
				log.Errorf("Error handling response: %v", err)
				// TODO: handle error properly
				panic(err)
			}
		}
	}()

	return c, nil
}

func (c *ClientTransport) Close() error {
	return c.conn.Close()
}

func (c *ClientTransport) HandleResponse(respBody []byte) error {
	// Decode the response
	response, err := c.Codec.Decode(respBody)
	if err != nil {
		return err
	}

	hopId, err := extractHopId(response)
	if err != nil {
		return err
	}

	// Find the corresponding request
	c.mu.Lock()
	ch, exists := c.msgCache[hopId]
	if exists {
		delete(c.msgCache, hopId)
	}
	c.mu.Unlock()

	if !exists {
		return errors.New("unexpected response hopId")
	}

	// Send the response to the waiting channel
	ch <- response
	close(ch)

	return nil
}

func (c *ClientTransport) SendMessage(request *mdd.Containers) (*mdd.Containers, error) {

	reqBody, err := c.Codec.Encode(request)
	if err != nil {
		return nil, err
	}

	hopId, err := extractHopId(request)
	if err != nil {
		return nil, err
	}

	ch := make(chan *mdd.Containers, 1)

	c.mu.Lock()
	c.msgCache[hopId] = ch
	c.mu.Unlock()

	err = Write(c.conn, reqBody)
	if err != nil {
		return nil, err
	}

	// Wait for the response from channel
	response := <-ch

	return response, nil
}

// need this in future when we start sending messages asynchronously
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
