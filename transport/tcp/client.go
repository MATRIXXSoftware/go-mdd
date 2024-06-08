package tcp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/mdd/field"
	"github.com/matrixxsoftware/go-mdd/transport"
)

type ClientTransport struct {
	conn  net.Conn
	Codec mdd.Codec
}

func NewClientTransport(addr string, codec mdd.Codec) (*ClientTransport, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &ClientTransport{
		conn:  conn,
		Codec: codec,
	}, nil
}

func (c *ClientTransport) Close() error {
	return c.conn.Close()
}

func (c *ClientTransport) SendMessage(request *mdd.Containers) (*mdd.Containers, error) {

	reqBody, err := c.Codec.Encode(request)
	if err != nil {
		return nil, err
	}

	respBody, err := c.send(reqBody)
	if err != nil {
		return nil, err
	}

	response, err := c.Codec.Decode(respBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientTransport) send(request []byte) ([]byte, error) {

	err := Write(c.conn, request)
	if err != nil {
		return nil, err
	}

	// Hard code 3 second for now. Make it configurable later
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	type Result struct {
		response []byte
		err      error
	}
	ch := make(chan Result, 1)

	go func() {
		response, err := Read(c.conn)
		ch <- Result{response, err}
	}()

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return nil, transport.ErrTimeout
		}
		return nil, ctx.Err()
	case res := <-ch:
		response := res.response
		err := res.err
		if err != nil {
			return nil, err
		}
		return response, nil
	}
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
