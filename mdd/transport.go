package mdd

type Client interface {
	SendMessage(request *Containers) (*Containers, error)
	Close() error
}

type Server interface {
	Listen() error
	Handler(handler func(*Containers) *Containers)
	Close() error
}
