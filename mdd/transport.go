package mdd

type ClientTransport interface {
	Send([]byte) ([]byte, error)
}

type ServerTransport interface {
	Handler(handler func([]byte) []byte)
}

type Client struct {
	codec     Codec
	transport ClientTransport
}

func NewClient(codec Codec, transport ClientTransport) (*Client, error) {
	return &Client{
		codec:     codec,
		transport: transport,
	}, nil
}

func (c *Client) SendMessage(request *Containers) (*Containers, error) {

	reqBody, err := c.codec.Encode(request)
	if err != nil {
		return nil, err
	}

	respBody, err := c.transport.Send(reqBody)
	if err != nil {
		return nil, err
	}

	response, err := c.codec.Decode(respBody)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type Server struct {
	codec     Codec
	transport ServerTransport
}

func NewServer(codec Codec, transport ServerTransport) (*Server, error) {
	return &Server{
		codec:     codec,
		transport: transport,
	}, nil
}

func (s *Server) Handler(handler func(*Containers) *Containers) {

	h := func(reqBody []byte) []byte {
		request, err := s.codec.Decode(reqBody)
		if err != nil {
			panic(err) // TODO
		}

		response := handler(request)

		respBody, err := s.codec.Encode(response)
		if err != nil {
			panic(err) // TODO
		}

		return respBody
	}

	s.transport.Handler(h)
}
