package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/transport/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ServerTransport struct {
	address string
	handler func(*mdd.Containers) (*mdd.Containers, error)
	Codec   mdd.Codec
	Tls     server.TLS
}

func NewServerTransport(addr string, codec mdd.Codec, opts ...server.Option) (*ServerTransport, error) {
	options := server.DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &ServerTransport{
		address: addr,
		Codec:   codec,
		Tls:     options.Tls,
	}, nil
}

func (s *ServerTransport) Listen() error {
	h2s := &http2.Server{}
	server := http.Server{
		Addr:    s.address,
		Handler: h2c.NewHandler(http.HandlerFunc(s.requestHandler), h2s),
	}

	if s.Tls.Enabled {
		return server.ListenAndServeTLS(s.Tls.CertFile, s.Tls.KeyFile)
	} else {
		return server.ListenAndServe()
	}
}

func (s *ServerTransport) Close() error {
	return nil
}

func (s *ServerTransport) Handler(handler func(*mdd.Containers) (*mdd.Containers, error)) {
	s.handler = handler
}

func (s *ServerTransport) requestHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	respBody, err := s.processMessage(reqBody)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, string(respBody))
}

func (s *ServerTransport) processMessage(reqBody []byte) ([]byte, error) {
	req, err := s.Codec.Decode(reqBody)
	if err != nil {
		return nil, err
	}

	response, err := s.handler(req)
	if err != nil {
		return nil, err
	}

	respBody, err := s.Codec.Encode(response)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}
