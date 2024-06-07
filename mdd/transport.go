package mdd

// These classes become redundant after the refactoring

type ClientTransport interface {
	SendMessage(*Containers) (*Containers, error)
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

func (c *Client) SendMessage(request *Containers) (*Containers, error) {
	return c.Transport.SendMessage(request)
}

type Server struct {
	Transport ServerTransport
}

func (s *Server) MessageHandler(handler func(*Containers) (*Containers, error)) {
	s.Transport.Handler(handler)
}
