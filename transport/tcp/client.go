package tcp

import (
	"net"
)

type ClientTransport struct {
	conn net.Conn
}

func NewClientTransport(addr string) (*ClientTransport, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &ClientTransport{
		conn: conn,
	}, nil
}

func (c *ClientTransport) Close() error {
	return c.conn.Close()
}

func (c *ClientTransport) Send(request []byte) ([]byte, error) {

	err := Encode(c.conn, request)
	if err != nil {
		return nil, err
	}

	response, err := Decode(c.conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}
