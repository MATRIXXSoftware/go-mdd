package http

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/transport/protocol/util"
	"github.com/matrixxsoftware/go-mdd/transport/server"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ServerTransport struct {
	address string
	handler func(*mdd.Containers) (*mdd.Containers, error)
	Codec   mdd.Codec
	tls     server.TLS
}

func NewServerTransport(addr string, codec mdd.Codec, opts ...server.Option) (*ServerTransport, error) {
	options := server.DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &ServerTransport{
		address: addr,
		Codec:   codec,
		tls:     options.Tls,
	}, nil
}

func (s *ServerTransport) Listen() error {
	h2s := &http2.Server{}
	server := http.Server{
		Addr:    s.address,
		Handler: h2c.NewHandler(http.HandlerFunc(s.requestHandler), h2s),
	}

	if s.tls.Enabled {
		cert, err := s.getCert()
		if err != nil {
			return err
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		server.TLSConfig = tlsConfig
		return server.ListenAndServeTLS("", "")
	} else {
		return server.ListenAndServe()
	}
}

func (s *ServerTransport) getCert() (tls.Certificate, error) {
	var err error
	var cert tls.Certificate
	var certPEM, keyPEM []byte
	if s.tls.SelfSignedCert {
		certPEM, keyPEM, err = util.GenerateSelfSignedCert()
	} else {
		certPEM, keyPEM, err = util.ReadCertAndKeyFiles(s.tls.CertFile, s.tls.KeyFile)
	}
	if err != nil {
		return cert, err
	}
	return tls.X509KeyPair(certPEM, keyPEM)
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
