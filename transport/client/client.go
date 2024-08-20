package client

import (
	"context"

	"github.com/matrixxsoftware/go-mdd/mdd"
)

type TLS struct {
	Enabled            bool
	InsecureSkipVerify bool
	CertFile           string
}

type Options struct {
	Tls TLS
}

type Option func(*Options)

func DefaultOptions() Options {
	return Options{
		Tls: TLS{
			Enabled:            false,
			InsecureSkipVerify: false,
			CertFile:           "",
		},
	}
}

func WithTls(tls TLS) func(*Options) {
	return func(o *Options) {
		o.Tls = tls
	}
}

type Transport interface {
	SendMessage(context.Context, *mdd.Containers) (*mdd.Containers, error)
	Close() error
}

type Client struct {
	Transport Transport
}

func (c *Client) SendMessage(ctx context.Context, request *mdd.Containers) (*mdd.Containers, error) {
	return c.Transport.SendMessage(ctx, request)
}
