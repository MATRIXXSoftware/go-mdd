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

	requestTransport := Transport{
		Containers: request,
		Codec:      c.codec,
	}
	err := requestTransport.Encode(c.conn)
	if err != nil {
		return nil, err
	}

	responseTransport := Transport{
		Containers: &Containers{},
		Codec:      c.codec,
	}

	err = responseTransport.Decode(c.conn)
	if err != nil {
		return nil, err
	}

	return responseTransport.Containers, nil
}
