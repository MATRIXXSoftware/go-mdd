package server

import "github.com/matrixxsoftware/go-mdd/mdd"

type TLS struct {
	Enable         bool
	SelfSignedCert bool
	CertFile       string
	KeyFile        string
}

type Options struct {
	Tls TLS
}

type Option func(*Options)

func DefaultOptions() Options {
	return Options{
		Tls: TLS{
			Enable:         false,
			SelfSignedCert: false,
			CertFile:       "",
			KeyFile:        "",
		},
	}
}

func WithTls(tls TLS) func(*Options) {
	return func(o *Options) {
		o.Tls = tls
	}
}

type Transport interface {
	Listen() error
	Handler(handler func(*mdd.Containers) (*mdd.Containers, error))
	Close() error
}

type Server struct {
	Transport Transport
}

func (s *Server) MessageHandler(handler func(*mdd.Containers) (*mdd.Containers, error)) {
	s.Transport.Handler(handler)
}
