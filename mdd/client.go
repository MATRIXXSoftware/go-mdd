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

func (c *Client) SendMessage(request *Containers) (*Containers, error) {
	err := request.Encode(c.conn)
	if err != nil {
		return nil, err
	}

	response := &Containers{}
	err = response.Decode(c.conn)
	if err != nil {
		return nil, err
	}

	return response, nil
}
