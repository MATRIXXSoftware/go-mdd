package http

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ServerTransport struct {
	address string
	handler func([]byte) []byte
}

func NewServerTransport(addr string) (*ServerTransport, error) {

	return &ServerTransport{
		address: addr,
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

func (s *ServerTransport) Handler(handler func([]byte) []byte) {
	s.handler = handler
}

func (s *ServerTransport) requestHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	respBody := s.handler(reqBody)

	fmt.Fprint(w, string(respBody))
}
