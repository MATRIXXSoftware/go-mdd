package tcp

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"os"
	"sync"

	"github.com/matrixxsoftware/go-mdd/mdd"
	"github.com/matrixxsoftware/go-mdd/transport"
	log "github.com/sirupsen/logrus"
)

type ServerTransport struct {
	ln      net.Listener
	handler func(*mdd.Containers) (*mdd.Containers, error)
	Codec   mdd.Codec
	mu      sync.Mutex
}

func NewServerTransport(addr string, codec mdd.Codec, opts ...transport.ServerOption) (*ServerTransport, error) {

	options := transport.DefaultServerOptions()
	for _, opt := range opts {
		opt(&options)
	}

	var ln net.Listener
	var err error

	tlsOptions := options.Tls
	if tlsOptions.Enable {
		var cert tls.Certificate
		var certPEM, keyPEM []byte
		if tlsOptions.SelfSignedCert {
			certPEM, keyPEM, err = generateSelfSignedCert()
		} else {
			certPEM, err = os.ReadFile(tlsOptions.CertFile)
			if err != nil {
				return nil, err
			}
			keyPEM, err = os.ReadFile(tlsOptions.KeyFile)
			if err != nil {
				return nil, err
			}
		}

		cert, err = tls.X509KeyPair(certPEM, keyPEM)
		if err != nil {
			return nil, err
		}

		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		ln, err = tls.Listen("tcp", addr, config)
	} else {
		ln, err = net.Listen("tcp", addr)
	}

	if err != nil {
		return nil, err
	}

	return &ServerTransport{
		ln:    ln,
		Codec: codec,
	}, nil
}

func (s *ServerTransport) Listen() error {
	for {
		// Accept a connection
		conn, err := s.ln.Accept()
		if err != nil {
			opErr, ok := err.(*net.OpError)
			if ok && opErr.Op == "accept" {
				log.Infof("ServerTransport shutting down")
				return nil
			}
			return err
		}
		// Spawn a new Goroutine for each incoming connection
		go s.handleConnection(conn)
	}
}

func (s *ServerTransport) Close() error {
	return s.ln.Close()
}

func (s *ServerTransport) Handler(handler func(*mdd.Containers) (*mdd.Containers, error)) {
	s.handler = handler
}

func (s *ServerTransport) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		reqBody, err := Read(conn)
		if err != nil {
			if err == io.EOF {
				log.Infof("%s Connection closed", connStr(conn))
			} else if err == io.ErrUnexpectedEOF {
				log.Errorf("%s Connection closed unexpectedly", connStr(conn))
			} else if err == context.DeadlineExceeded {
				log.Errorf("%s Timeout reading from connection", connStr(conn))
			} else {
				log.Errorf("%s Error reading from connection: %s", connStr(conn), err)
			}
			return
		}

		if reqBody == nil {
			continue
		}

		go func() {
			respBody, err := s.processRequest(reqBody)
			if err != nil {
				log.Errorf("%s %s", connStr(conn), err)
				return
			}

			// Multiple goroutines can write to the same connection
			// Therefore we need to lock the write operation here
			s.mu.Lock()
			err = Write(conn, respBody)
			s.mu.Unlock()

			if err != nil {
				log.Errorf("%s %s", connStr(conn), err)
				return
			}
		}()
	}
}

func (s *ServerTransport) processRequest(reqBody []byte) ([]byte, error) {
	req, err := s.Codec.Decode(reqBody)
	if err != nil {
		return nil, err
	}

	missingHopId := false
	hopId, err := extractHopId(req)
	if err != nil {
		log.Debugf("Error extracting hopId: %s", err)
		missingHopId = true
	}

	response, err := s.handler(req)
	if err != nil {
		return nil, err
	}

	if !missingHopId {
		err := injectHopId(response, hopId)
		if err != nil {
			return nil, err
		}
	}

	respBody, err := s.Codec.Encode(response)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}
