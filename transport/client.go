package transport

import "net"

type MddClient struct {
	conn net.Conn
}

func NewClient(addr string) (*MddClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &MddClient{conn: conn}, nil
}

func (c *MddClient) Close() error {
	return c.conn.Close()
}

func (c *MddClient) SendMessage(msg *ExampleMessage) (*ExampleMessage, error) {
	err := msg.Encode(c.conn)
	if err != nil {
		return nil, err
	}

	response := &ExampleMessage{}
	err = response.Decode(c.conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}
