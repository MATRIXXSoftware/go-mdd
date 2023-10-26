package mdd

import "net"

type Client struct {
	conn  net.Conn
	codec Codec
}

func NewClient(addr string, codec Codec) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:  conn,
		codec: codec,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) SendMessage(request *Containers) (*Containers, error) {

	err := Encode(c.conn, c.codec, request)
	if err != nil {
		return nil, err
	}

	response, err := Decode(c.conn, c.codec)
	if err != nil {
		return nil, err
	}

	return response, nil
}
