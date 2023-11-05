package mdd

type ClientTransport interface {
	Send([]byte) ([]byte, error)
	Close() error
}

type ServerTransport interface {
	Listen() error
	Handler(handler func([]byte) ([]byte, error))
	Close() error
}

type Client struct {
	Codec     Codec
	Transport ClientTransport
}

func (c *Client) SendMessage(request *Containers) (*Containers, error) {

	reqBody, err := c.Codec.Encode(request)
	if err != nil {
		return nil, err
	}

	respBody, err := c.Transport.Send(reqBody)
	if err != nil {
		return nil, err
	}

	response, err := c.Codec.Decode(respBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type Server struct {
	Codec     Codec
	Transport ServerTransport
}

func (s *Server) MessageHandler(handler func(*Containers) *Containers) {

	h := func(reqBody []byte) ([]byte, error) {
		request, err := s.Codec.Decode(reqBody)
		if err != nil {
			return nil, err
		}

		response := handler(request)

		respBody, err := s.Codec.Encode(response)
		if err != nil {
			return nil, err
		}

		return respBody, nil
	}

	s.Transport.Handler(h)
}
