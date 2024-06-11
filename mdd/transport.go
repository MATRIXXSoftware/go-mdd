package mdd

import "context"

// These classes become redundant after the refactoring

type ClientTransport interface {
	SendMessage(context.Context, *Containers) (*Containers, error)
	Close() error
}

type ServerTransport interface {
	Listen() error
	Handler(handler func(*Containers) (*Containers, error))
	Close() error
}

type Client struct {
	Transport ClientTransport
}

func (c *Client) SendMessage(ctx context.Context, request *Containers) (*Containers, error) {
	return c.Transport.SendMessage(ctx, request)
}

type Server struct {
	Transport ServerTransport
}

func (s *Server) MessageHandler(handler func(*Containers) (*Containers, error)) {
	s.Transport.Handler(handler)
}
