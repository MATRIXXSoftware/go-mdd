package http

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	log "github.com/sirupsen/logrus"
)

type ServerTransport struct {
	address string
	handler func([]byte) ([]byte, error)
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

func (s *ServerTransport) Handler(handler func([]byte) ([]byte, error)) {
	s.handler = handler
}

func (s *ServerTransport) requestHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("received request: %v", r)
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	log.Debugf("received request body: %s", reqBody)
	respBody, err := s.handler(reqBody)
	if err != nil {
		panic(err)
	}

	log.Debugf("response body: %s", string(respBody))
	fmt.Fprint(w, string(respBody))
}
