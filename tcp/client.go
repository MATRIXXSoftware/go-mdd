package tcp

import (
	"net"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

type Client struct {
	conn  net.Conn
	codec mdd.Codec
}

func NewClient(addr string, codec mdd.Codec) (*Client, error) {
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

func (c *Client) SendMessage(request *mdd.Containers) (*mdd.Containers, error) {

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
