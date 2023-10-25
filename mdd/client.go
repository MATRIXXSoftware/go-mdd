package mdd

import "net"

type Client struct {
	conn net.Conn
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendMessage(msg *ExampleMessage) (*ExampleMessage, error) {
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
