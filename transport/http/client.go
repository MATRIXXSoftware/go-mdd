package http

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"

	"golang.org/x/net/http2"
)

type ClientTransport struct {
	httpClient http.Client
	address    string
}

func NewClientTransport(addr string) (*ClientTransport, error) {
	httpClient := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}

	return &ClientTransport{
		httpClient: httpClient,
		address:    addr,
	}, nil
}

func (c *ClientTransport) Close() error {
	return nil
}

func (c *ClientTransport) Send(reqBody []byte) ([]byte, error) {

	req, err := http.NewRequest("POST", "http://"+c.address, bytes.NewReader(reqBody))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil || len(respBody) == 0 {
		return nil, err
	}

	return respBody, nil
}
