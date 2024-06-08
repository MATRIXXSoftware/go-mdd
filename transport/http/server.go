package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ServerTransport struct {
	address string
	handler func(*mdd.Containers) (*mdd.Containers, error)
	Codec   mdd.Codec
}

func NewServerTransport(addr string, codec mdd.Codec) (*ServerTransport, error) {

	return &ServerTransport{
		address: addr,
		Codec:   codec,
	}, nil
}

func (s *ServerTransport) Listen() error {
	h2s := &http2.Server{}
	server := http.Server{
		Addr:    s.address,
		Handler: h2c.NewHandler(http.HandlerFunc(s.requestHandler), h2s),
	}
	return server.ListenAndServe()
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
