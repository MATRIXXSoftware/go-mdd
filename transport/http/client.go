package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"golang.org/x/net/http2"
)

type ClientTransport struct {
	httpClient http.Client
	address    string
	Codec      mdd.Codec
}

func NewClientTransport(addr string, codec mdd.Codec) (*ClientTransport, error) {
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
		Codec:      codec,
	}, nil
}

func (c *ClientTransport) Close() error {
	return nil
}

func (c *ClientTransport) SendMessage(ctx context.Context, request *mdd.Containers) (*mdd.Containers, error) {
	reqBody, err := c.Codec.Encode(request)
	if err != nil {
		return nil, err
	}

	respBody, err := c.send(ctx, reqBody)
	if err != nil {
		return nil, err
	}

	response, err := c.Codec.Decode(respBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *ClientTransport) send(ctx context.Context, reqBody []byte) ([]byte, error) {

	req, err := http.NewRequestWithContext(ctx, "POST", "http://"+c.address, bytes.NewReader(reqBody))

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
